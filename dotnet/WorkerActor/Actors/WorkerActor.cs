using Dapr.Actors;
using Dapr.Actors.Runtime;
using Dapr.Client;


namespace BasicActorSamples.Actors
{
    public interface IWorkerActor : IActor
    {
        Task<string> GetState();
        Task SetState(string state);

        Task<bool> DeleteOrder(string orderId);
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
    }
}