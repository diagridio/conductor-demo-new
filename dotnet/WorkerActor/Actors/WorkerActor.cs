using Dapr.Actors;
using Dapr.Actors.Runtime;
using Dapr.Client;
using System.Threading.Tasks;


namespace BasicActorSamples.Actors
{
    public interface IWorkerActor : IActor
    {
        Task<string> GetState();
        Task SetState(string state);

        Task<bool> DeleteOrder(string orderId);

        Task<string> SetDataAsync(MyData data);
        Task<MyData> GetDataAsync();
        Task RegisterReminder();
        Task UnregisterReminder();
        Task<IActorReminder> GetReminder();
        Task RegisterTimer();
        Task UnregisterTimer();
    }

    public class MyData
    {
        public string PropertyA { get; set; }
        public string PropertyB { get; set; }

        public override string ToString()
        {
            var propAValue = this.PropertyA == null ? "null" : this.PropertyA;
            var propBValue = this.PropertyB == null ? "null" : this.PropertyB;
            return $"PropertyA: {propAValue}, PropertyB: {propBValue}";
        }
    }
    
    public class WorkerActor : Actor, IWorkerActor
    {
        string STATESTORE = "oms.state.makeline";
        public WorkerActor(ActorHost host) : base(host)
        {
        }
        public async Task SetState(string state)
        {
            using var client = new DaprClientBuilder().Build();

            //Using Dapr SDK to save and get state
            await client.SaveStateAsync(STATESTORE, "order_1", state);
        }

        public async Task<string> GetState()
        {
            using var client = new DaprClientBuilder().Build();
            return await client.GetStateAsync<string>(STATESTORE, "order_1");
        }

         public async Task<bool> DeleteOrder(string orderId)
        {
            try{
                using var client = new DaprClientBuilder().Build();

                CancellationTokenSource source = new CancellationTokenSource();
                CancellationToken cancellationToken = source.Token;
                
                await client.DeleteStateAsync(STATESTORE, orderId, cancellationToken: cancellationToken);        

                return true;
            }catch(Exception ex){
                return false;
            }
            
        }

        // Mock methods to be displayed on Conductor

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

        /// <summary>
        /// Register MyReminder reminder with the actor
        /// </summary>
        public async Task RegisterReminder()
        {
            await this.RegisterReminderAsync(
                "MyReminder",              // The name of the reminder
                null,                      // User state passed to IRemindable.ReceiveReminderAsync()
                TimeSpan.FromSeconds(1),   // Time to delay before invoking the reminder for the first time
                TimeSpan.FromSeconds(60));  // Time interval between reminder invocations after the first invocation
        }

        // /// <summary>
        // /// Get MyReminder reminder details with the actor
        // /// </summary>
        public async Task<IActorReminder> GetReminder()
        {
            return await this.GetReminderAsync("MyReminder");
        }

        /// <summary>
        /// Unregister MyReminder reminder with the actor
        /// </summary>
        public Task UnregisterReminder()
        {
            Console.WriteLine("Unregistering My Reminder...");
            return this.UnregisterReminderAsync("MyReminder");
        }

        // <summary>
        // Implement IRemindeable.ReceiveReminderAsync() which is call back invoked when an actor reminder is triggered.
        // </summary>
        public Task ReceiveReminderAsync(string reminderName, byte[] state, TimeSpan dueTime, TimeSpan period)
        {
            Console.WriteLine("ReceiveReminderAsync is called!");
            return Task.CompletedTask;
        }

        /// <summary>
        /// Register MyTimer timer with the actor
        /// </summary>
        public Task RegisterTimer()
        {
            return this.RegisterTimerAsync(
                "MyTimer",                  // The name of the timer
                nameof(this.OnTimerCallBack),       // Timer callback
                null,                       // User state passed to OnTimerCallback()
                TimeSpan.FromSeconds(1),    // Time to delay before the async callback is first invoked
                TimeSpan.FromSeconds(30));   // Time interval between invocations of the async callback
        }

        /// <summary>
        /// Unregister MyTimer timer with the actor
        /// </summary>
        public Task UnregisterTimer()
        {
            Console.WriteLine("Unregistering MyTimer...");
            return this.UnregisterTimerAsync("MyTimer");
        }

        /// <summary>
        /// Timer callback once timer is expired
        /// </summary>
        private Task OnTimerCallBack(byte[] data)
        {
            Console.WriteLine("OnTimerCallBack is called!");
            this.UnregisterTimer();
            this.UnregisterReminder();
            return Task.CompletedTask;
        }
    }
}