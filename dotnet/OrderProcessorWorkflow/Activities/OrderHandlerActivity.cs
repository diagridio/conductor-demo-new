using Dapr.Client;
using Dapr.Workflow;
using OrderProcessorWorkflow.Models;
using Dapr.Actors;
using Dapr.Actors.Client;
using WorkerActor.Interface;

namespace OrderProcessorWorkflow.Activities
{
    public class OrderHandlerActivity : WorkflowActivity<OrderInput, object>
    {
        readonly ILogger _logger;
        readonly DaprClient _client;

        private const string LoyaltyServiceId = "loyalty-service";
        private const string ReceiptServiceId = "receipt-generation-service";
        private const string MakeLineServiceId = "make-line-service";
        private const string WorkerActorType = "WorkerActor";

        public OrderHandlerActivity(ILoggerFactory loggerFactory, DaprClient client)
        {
            _logger = loggerFactory.CreateLogger<OrderHandlerActivity>();
            _client = client;
        }

        public override async Task<object> RunAsync(WorkflowActivityContext context, OrderInput req)
        {
            
            // Checks for activity type and calls the appropriate method
            if (req.requestType == OrderRequestType.Loyalty) {
                return await HandleLoyaltyRequest(req.orderSummary);
            }
            else if (req.requestType == OrderRequestType.Receipt) {
                return await HandleReceiptRequest(req.orderSummary);
            } 
            else if (req.requestType == OrderRequestType.MakeLine) {
                return await HandleMakeLineRequest(req.orderSummary);
            }
            else if (req.requestType == OrderRequestType.Complete) {
                return await HandleCompleteOrderRequest(req.orderSummary);
            }
            else{
                return new OrderResult(Processed: false);
            }
        }

        private async Task<OrderResult> HandleLoyaltyRequest(OrderSummary orderSummary){
            _logger.LogInformation("Processing loyalty points {orderId}.",orderSummary.OrderId);

            var request = _client.CreateInvokeMethodRequest<Object>(LoyaltyServiceId, "loyalty", orderSummary);
            var response = await _client.InvokeMethodWithResponseAsync(request);

                if (!response.IsSuccessStatusCode)
                {
                    _logger.LogInformation("Loyalty points update was unsuccessful: {0} {1} {2}", (int)response.StatusCode, response.StatusCode, await response.Content.ReadAsStringAsync());
                    return new OrderResult(Processed: false);
                }
                else
                {
                    _logger.LogInformation("Loyalty updaed for {3}: {0} {1}.", orderSummary.FirstName, orderSummary.LastName, orderSummary.LoyaltyId);
                    return new OrderResult(Processed:true);
                    
                }
        }

         private async Task<OrderResult> HandleReceiptRequest(OrderSummary orderSummary){
            _logger.LogInformation(
                "Generating receipt for {orderId}.",
                orderSummary.OrderId);

            var request = _client.CreateInvokeMethodRequest<Object>(ReceiptServiceId, "receipt", orderSummary);
            var response = await _client.InvokeMethodWithResponseAsync(request);

                if (!response.IsSuccessStatusCode)
                {
                    _logger.LogInformation("Receipt generations was unsuccessful: {0} {1} {2}", (int)response.StatusCode, response.StatusCode, await response.Content.ReadAsStringAsync());
                    return new OrderResult(Processed: false);
                }
                else
                {
                    _logger.LogInformation("Receipt generated for customer {0} {1}.", orderSummary.FirstName, orderSummary.LastName);
                    return new OrderResult(Processed: true);
                    
                }
        }

        private async Task<OrderResult> HandleMakeLineRequest(OrderSummary orderSummary){
            _logger.LogInformation(
                "Starting make-line process for {orderId}.",
                orderSummary.OrderId);

            var request = _client.CreateInvokeMethodRequest<Object>(MakeLineServiceId, "makeline", orderSummary);
            var response = await _client.InvokeMethodWithResponseAsync(request);

                if (!response.IsSuccessStatusCode)
                {
                    _logger.LogInformation("Makeline process failed: {0} {1} {2}", (int)response.StatusCode, response.StatusCode, await response.Content.ReadAsStringAsync());
                    return new OrderResult(Processed: false);
                }
                else
                {
                    _logger.LogInformation("Make-line process completed succesfully for order {0}.", orderSummary.OrderId);
                    return new OrderResult(Processed: true);
                    
                }
        }

        private async Task<OrderResult> HandleCompleteOrderRequest(OrderSummary orderSummary){
            _logger.LogInformation(
                "Starting order completion for {orderId}.",
                orderSummary.OrderId);

            var oId = orderSummary.OrderId.ToString();

            // In the Client Application
            //var actorId = new ActorId(orderSummary.OrderId.ToString());

            var proxy = ActorProxy.Create<IWorkerActor>(ActorId.CreateRandom(), WorkerActorType);

            var response = await proxy.CompleteOrder(oId);

                if (!response)
                {
                    _logger.LogInformation("Worker Actor failed to complete order: {0}", orderSummary.OrderId);
                    return new OrderResult(Processed: false);
                }
                else
                {
                    _logger.LogInformation("Worker Actor completed order {0}.", orderSummary.OrderId);
                    return new OrderResult(Processed: true);
                    
                }
        }
    }
}
