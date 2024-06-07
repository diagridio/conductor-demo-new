using Dapr.Actors;

namespace BasicActorSamples.Actors
{
    public interface IWorkerActor : IActor
    {
        Task<string> GetState();
        Task SetState(string state);

        Task<bool> DeleteOrder(string orderId);
    }
}
