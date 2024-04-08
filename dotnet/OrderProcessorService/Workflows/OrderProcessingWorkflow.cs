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
                // End the workflow here since we can;t create the receipt
                await context.CallActivityAsync(
                    nameof(NotifyActivity),
                    new Notification($"Failed to generate receipt to {order.FirstName} {order.LastName}"));

                context.SetCustomStatus("Stopped order process due to receipt generation issues.");

                return new OrderResult(Processed: false);
            }

            // Update loyalty points for the customer
              await context.CallActivityAsync<OrderResult>(
                nameof(OrderHandlerActivity),
                new OrderInput(OrderRequestType.Loyalty, order));
                //this.defaultActivityRetryOptions);
            

            if (!result.Processed)
            {
                // End the workflow here we can't update loyalty points for the customer
                await context.CallActivityAsync(
                    nameof(NotifyActivity),
                    new Notification($"Failed to updated loyalty points for {order.FirstName} {order.LastName}"));
                    
                context.SetCustomStatus("Stopped order process due to failure to update the loyalty points.");
                
                return new OrderResult(Processed: false);
            }
            
            // Determine if the make-line process is successful
              await context.CallActivityAsync<OrderResult>(
                nameof(OrderHandlerActivity),
                new OrderInput(OrderRequestType.MakeLine, order));
                //this.defaultActivityRetryOptions);
            
            // If there the make line fails, stop and let the user know
            if (!result.Processed)
            {
                // End the workflow here due to make-line issues
                await context.CallActivityAsync(
                    nameof(NotifyActivity),
                    new Notification($"Failed to complete make-line process for order {order.OrderId}"));

                    context.SetCustomStatus("Stopped order process due to failure to complete make-line process.");

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
