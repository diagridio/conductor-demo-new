using Dapr.Actors;

namespace WorkerActor.Interface
{
    public interface IWorkerActor : IActor
    {   
        Task<bool> CompleteOrder(string orderId);
    }
}
