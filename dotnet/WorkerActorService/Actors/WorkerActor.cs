using Dapr.Actors;
using Dapr.Actors.Runtime;
using System;
using WorkerActor.Interface;
using System.Threading.Tasks;
using Dapr.Client;


namespace WorkerActorService
{
    internal class WorkerActor : Actor, IWorkerActor
    {

        private readonly string storeName = "oms.state.makeline"; // Replace with your state store name

        private DaprClient client;

        // The constructor must accept ActorHost as a parameter, and can also accept additional
        // parameters that will be retrieved from the dependency injection container
        //
        /// <summary>
        /// Initializes a new instance of MyActor
        /// </summary>
        /// <param name="host">The Dapr.Actors.Runtime.ActorHost that will host this actor instance.</param>
        public WorkerActor(ActorHost host)
            : base(host)
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

        }
    

        /// <summary>
        /// This method is called whenever an actor is activated.
        /// An actor is activated the first time any of its methods are invoked.
        /// </summary>
        protected override Task OnActivateAsync()
        {
            // Provides opportunity to perform some optional setup.
            Console.WriteLine($"Activating actor id: {this.Id}");
            return Task.CompletedTask;
        }

        /// <summary>
        /// This method is called whenever an actor is deactivated after a period of inactivity.
        /// </summary>
        protected override Task OnDeactivateAsync()
        {
            // Provides Opporunity to perform optional cleanup.
            Console.WriteLine($"Deactivating actor id: {this.Id}");
            return Task.CompletedTask;
        }

    }
}
