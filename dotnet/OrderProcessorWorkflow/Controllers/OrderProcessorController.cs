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
        private readonly Random _random;
        //string WorkflowName = "OrderProcessingWorkflow";
        //private readonly StateOptions _stateOptions = new StateOptions(){ Concurrency = ConcurrencyMode.FirstWrite, Consistency = ConsistencyMode.Eventual };

        public OrderProcessorController(ILogger<OrderProcessorController> logger, DaprWorkflowClient daprClient)
        {
            _logger = logger;
            _daprClient = daprClient;
            _random = new Random();
        }
        
        [Dapr.Topic(PubSubName, OrderTopic)]
        [HttpPost("/orders")]
        public async Task<IActionResult> Process([FromBody] OrderSummary orderSummary)
        {   
            if (orderSummary is not null) 
            {
                var orderId = orderSummary.OrderId.ToString();
                var instanceId = _random.Next(10000).ToString();

                _logger.LogInformation("Received Order. Initializing workflow {instanceId}.", instanceId);

                    // Start the workflow using the order ID as the workflow ID
                _logger.LogInformation("Starting order {orderId}", orderId);
                    await _daprClient.ScheduleNewWorkflowAsync(
                    name: nameof(OrderProcessingWorkflow),
                    input: orderSummary,
                    instanceId: instanceId);

                // Wait for the workflow to start and confirm the input
                WorkflowState state = await _daprClient.WaitForWorkflowStartAsync(instanceId: instanceId);

                // Wait for the workflow to complete
                while (true)
                {
                    using var cts = new CancellationTokenSource(TimeSpan.FromSeconds(30));
                    try
                    {
                        state = await _daprClient.WaitForWorkflowCompletionAsync(
                            instanceId: instanceId,
                            cancellation: cts.Token);
                        break;
                    }
                    catch (OperationCanceledException)
                    {
                        _logger.LogInformation("Waiting for {orderId} to complete.", orderId);
                        return BadRequest();
                    }
                }
                if (state.RuntimeStatus == WorkflowRuntimeStatus.Completed)
                {
                    //await _daprClient.PurgeInstanceAsync(orderId);

                    OrderResult result = state.ReadOutputAs<OrderResult>();
                    if (result.Processed)
                    {
                        _logger.LogInformation("Order workflow is {state.RuntimeStatus} and was processed successfully ({result}).", state.RuntimeStatus, result);
                        return Ok();
                    }
                    else
                    {
                        _logger.LogInformation("Order workflow is  {state.RuntimeStatus} but the order was not processed.", state.RuntimeStatus);
                        return BadRequest();
                    }
                }
                else if (state.RuntimeStatus == WorkflowRuntimeStatus.Failed)
                {
                    //await _daprClient.PurgeInstanceAsync(orderId);

                    _logger.LogInformation("The workflow failed - {state.FailureDetails}", state.FailureDetails);
                    return BadRequest();
                }

            }
            
            return BadRequest();
        }

    }
}