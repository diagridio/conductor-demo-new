using Dapr.Actors;
using Dapr.Actors.Runtime;
using System.Threading.Tasks;

namespace WorkerActor.Interface
{
    public interface IWorkerActor : IActor
    {       
        Task<string> SetDataAsync(MyData data);
        Task<MyData> GetDataAsync();
        Task<bool> CompleteOrder(string orderId);
    }

    public class MyData
    {
        public string PropertyA { get; set; }
        public string PropertyB { get; set; }

        public override string ToString()
        {
            var propAValue = this.PropertyA == null ? "null" : this.PropertyA;
            var propBValue = this.PropertyB == null ? "null" : this.PropertyB;
            return $"Proper tyA: {propAValue}, PropertyB: {propBValue}";
        }
    }
}
