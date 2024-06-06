using Dapr.Actors;
using Dapr.Actors.Runtime;
using Dapr.Client;

namespace WorkerActor.Actors
{
    public interface IWorker : IActor
    {
        Task<bool> CompleteOrder(string orderId);
    }
    
    [Actor(TypeName = "WorkerActor")]
    public class Worker : Actor, IWorker
    {
        private readonly string storeName = "oms.state.makeline"; // Replace with your state store name

        private DaprClient client;

        public Worker(ActorHost host) : base(host)
        {
            client = new DaprClientBuilder().Build();
        }

        public async Task<bool> CompleteOrder(string orderId)
        {
            try
            {
                await client.DeleteStateAsync(storeName, orderId);
                Console.WriteLine($"Order {orderId} deleted.");
                
                return true;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error deleting order: {ex.Message}");
                return false;
            }

            return true;
        }
    }
}