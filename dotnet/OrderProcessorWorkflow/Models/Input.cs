namespace OrderProcessorWorkflow.Models
{
    public record OrderInput(OrderRequestType requestType, OrderSummary orderSummary);
    public record OrderResult(bool Processed);
    public record Notification(string Message);

    public enum OrderRequestType { Receipt, Loyalty, MakeLine, Complete }

}