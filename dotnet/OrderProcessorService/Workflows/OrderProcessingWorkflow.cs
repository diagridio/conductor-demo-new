using System.Diagnostics;
using Dapr.Workflow;
using OrderProcessorService.Activities;
using OrderProcessorService.Models;

namespace OrderProcessorService.Workflows
{
    public class OrderProcessingWorkflow : Workflow<OrderSummary, OrderResult>
    {
        // readonly WorkflowTaskOptions defaultActivityRetryOptions = new WorkflowTaskOptions
        // {
        //     // NOTE: Beware that changing the number of retries is a breaking change for existing workflows.
        //     RetryPolicy = new WorkflowRetryPolicy(
        //         maxNumberOfAttempts: 3,
        //         firstRetryInterval: TimeSpan.FromSeconds(5)),
        // };

        public async override Task<OrderResult> RunAsync(WorkflowContext context, OrderSummary order)
        {
            string orderId = context.InstanceId;

             // Notify the user that an order has come through
            await context.CallActivityAsync(
                nameof(NotifyActivity),
                new Notification($"Received order {orderId} for {order.FirstName} {order.LastName} at ${order.StoreId}"));

          // Generating receipt for the order
            OrderResult result = await context.CallActivityAsync<OrderResult>(
                nameof(OrderHandlerActivity),
                new OrderInput(OrderRequestType.Receipt, order));
                //this.defaultActivityRetryOptions);

            // If receipt generation fails, let the user know 
            if (!result.Processed)
            {
                // End the workflow here since we don't have sufficient inventory
                await context.CallActivityAsync(
                    nameof(NotifyActivity),
                    new Notification($"Failed to generate receipt to {order.FirstName} {order.LastName}"));
                return new OrderResult(Processed: false);
            }

            // Determine if there is enough of the item available for purchase by checking the inventory
              await context.CallActivityAsync<OrderResult>(
                nameof(OrderHandlerActivity),
                new OrderInput(OrderRequestType.Loyalty, order));
                //this.defaultActivityRetryOptions);
            
            // If there is insufficient inventory, fail and let the user know 
            if (!result.Processed)
            {
                // End the workflow here since we don't have sufficient inventory
                await context.CallActivityAsync(
                    nameof(NotifyActivity),
                    new Notification($"Failed to updated loyalty points for {order.FirstName} {order.LastName}"));
                return new OrderResult(Processed: false);
            }

            // Let them know their order was processed
            await context.CallActivityAsync(
                nameof(NotifyActivity),
                new Notification($"Order {orderId} has completed!"));

            // End the workflow with a success result
            return new OrderResult(Processed: true);
        }


    }
}
