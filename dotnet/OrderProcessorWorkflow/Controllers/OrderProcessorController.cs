using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Threading.Tasks;
using Dapr;
using Dapr.Client;
using Microsoft.AspNetCore.Cors;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using OrderProcessorWorkflow.Models;
using System.Text.Json;
using System.Text;
using Microsoft.AspNetCore.Http.HttpResults;
using Dapr.Workflow;
using OrderProcessorWorkflow.Workflows;

namespace OrderProcessorWorkflow.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class OrderProcessorController : ControllerBase
    {
        private readonly ILogger<OrderProcessorController> _logger;
        private const string OrderTopic = "orders";
        private const string PubSubName = "oms.pubsub";
        private readonly DaprWorkflowClient _daprClient;
        //private readonly StateOptions _stateOptions = new StateOptions(){ Concurrency = ConcurrencyMode.FirstWrite, Consistency = ConsistencyMode.Eventual };

        public OrderProcessorController(ILogger<OrderProcessorController> logger, DaprWorkflowClient daprClient)
        {
            _logger = logger;
            _daprClient = daprClient;
        }

        [Dapr.Topic(PubSubName, OrderTopic)]
        [HttpPost("/orders")]
        public async Task<IActionResult> Process([FromBody] OrderSummary orderSummary)
        {
            if (orderSummary is null)
            {
                // Malformed message - there's nothing we can do with it, and returning an
                // error status here would just cause Dapr to redeliver the same bad payload
                // forever. Ack it and move on.
                _logger.LogWarning("Received a null/unparseable order summary; dropping message.");
                return Ok();
            }

            var orderId = orderSummary.OrderId.ToString();

            // Use the order ID - not a random number - as the workflow instance ID. Kafka
            // (and Dapr's pub/sub retry behavior) can redeliver the same "orders" message
            // more than once; a deterministic instance ID lets us detect that below and
            // avoid starting a second workflow run - and re-triggering loyalty/receipt/
            // make-line/actor processing - for an order that's already being handled.
            var instanceId = orderId;

            _logger.LogInformation("Received order {orderId}.", orderId);

            WorkflowState existingState = await _daprClient.GetWorkflowStateAsync(
                instanceId: instanceId,
                getInputsAndOutputs: false);

            if (existingState.Exists)
            {
                if (existingState.IsWorkflowCompleted)
                {
                    // Duplicate delivery of an order we've already finished (successfully or
                    // not). Ack it without reprocessing.
                    _logger.LogInformation(
                        "Order {orderId} was already processed (status: {status}); acknowledging duplicate delivery without reprocessing.",
                        orderId, existingState.RuntimeStatus);
                    return Ok();
                }

                _logger.LogInformation(
                    "Order {orderId} is already in progress (status: {status}); waiting for it to complete.",
                    orderId, existingState.RuntimeStatus);
            }
            else
            {
                _logger.LogInformation("Starting workflow for order {orderId}.", orderId);
                await _daprClient.ScheduleNewWorkflowAsync(
                    name: nameof(OrderProcessingWorkflow),
                    input: orderSummary,
                    instanceId: instanceId);

                // Wait for the workflow to start and confirm the input
                await _daprClient.WaitForWorkflowStartAsync(instanceId: instanceId);
            }

            WorkflowState state;
            using (var cts = new CancellationTokenSource(TimeSpan.FromSeconds(30)))
            {
                try
                {
                    state = await _daprClient.WaitForWorkflowCompletionAsync(
                        instanceId: instanceId,
                        cancellation: cts.Token);
                }
                catch (OperationCanceledException)
                {
                    // The workflow hasn't finished within our wait window, but it IS running.
                    // Ack the message instead of returning an error status: the workflow keeps
                    // running to completion on its own, and failing the delivery here would
                    // only cause Dapr to redeliver this Kafka message - which, thanks to the
                    // deterministic instance ID above, would just re-check on the same
                    // in-flight workflow rather than making anything happen faster.
                    _logger.LogInformation(
                        "Order {orderId} is still processing after 30s; acknowledging delivery and letting the workflow continue in the background.",
                        orderId);
                    return Ok();
                }
            }

            if (state is null)
            {
                _logger.LogError("Workflow state returned as null for order {orderId} (instance {instanceId}).", orderId, instanceId);
                return Ok();
            }

            if (state.RuntimeStatus == WorkflowRuntimeStatus.Completed)
            {
                OrderResult? result = state.ReadOutputAs<OrderResult>();
                if (result is not null && result.Processed)
                {
                    _logger.LogInformation("Order {orderId} completed successfully ({result}).", orderId, result);
                }
                else
                {
                    _logger.LogWarning("Order {orderId} workflow completed but reported Processed=false.", orderId);
                }
            }
            else if (state.RuntimeStatus == WorkflowRuntimeStatus.Failed)
            {
                _logger.LogWarning("Order {orderId} workflow failed: {failureDetails}", orderId, state.FailureDetails);
            }

            // Always ack the Kafka message once we've handled it, whether the order was
            // processed successfully or not: the deterministic instance ID above means a
            // redelivery would just no-op against this same completed/failed workflow
            // instead of reprocessing the order, so failing the delivery here has no upside -
            // only the downside of triggering retry storms that hammer the actor runtime
            // (which is what was happening before this fix).
            return Ok();
        }

    }
}
