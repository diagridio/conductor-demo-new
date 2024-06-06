using Dapr.Actors;
using Dapr.Actors.Runtime;
using System;
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

        /// <summary>
        /// Set MyData into actor's private state store
        /// </summary>
        /// <param name="data">the user-defined MyData which will be stored into state store as "my_data" state</param>
        public async Task<string> SetDataAsync(MyData data)
        {
            // Data is saved to configured state store implicitly after each method execution by Actor's runtime.
            // Data can also be saved explicitly by calling this.StateManager.SaveStateAsync();
            // State to be saved must be DataContract serializable.
            await this.StateManager.SetStateAsync<MyData>(
                "my_data",  // state name
                data);      // data saved for the named state "my_data"

            return "Success";
        }

        /// <summary>
        /// Get MyData from actor's private state store
        /// </summary>
        /// <return>the user-defined MyData which is stored into state store as "my_data" state</return>
        public Task<MyData> GetDataAsync()
        {
            // Gets state from the state store.
            return this.StateManager.GetStateAsync<MyData>("my_data");
        }

    }
}
