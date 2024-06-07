using Dapr.Actors;
using Dapr.Actors.Runtime;
using Dapr.Client;


namespace BasicActorSamples.Actors
{
    public interface IHelloWorld : IActor
    {
        Task<string> SayHelloWorld();
        Task SetState(string state);
        Task<string> GetState();
        Task<string> SayHello(string name);
    }
    
    public class HelloWorldActor : Actor, IHelloWorld
    {
        string STATESTORE = "oms.state.makeline";
        public HelloWorldActor(ActorHost host) : base(host)
        {
        }

        public Task<string> SayHelloWorld()
        {
            return Task.FromResult("Hello World!");
        }

        public Task<string> SayHello(string name)
        {
            return Task.FromResult($"Hello {name}!");
        }

        // public async Task SetState(string state)
        // {
        //    await StateManager.SetStateAsync(STATESTORE, state);
        // }

        // public async Task<string> GetState()
        // {
        //     return await StateManager.GetStateAsync<string>(STATESTORE);
        // }

        public async Task SetState(string state)
        {
            using var client = new DaprClientBuilder().Build();
            Random random = new Random();
            int orderId = random.Next(1,1000);
            //Using Dapr SDK to save and get state
            await client.SaveStateAsync(STATESTORE, "order_1", orderId.ToString());
        }

        public async Task<string> GetState()
        {
            using var client = new DaprClientBuilder().Build();
            return await client.GetStateAsync<string>(STATESTORE, "order_1");
        }
    }
}