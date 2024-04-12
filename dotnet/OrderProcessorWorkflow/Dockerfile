FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src
COPY . .
RUN dotnet restore "OrderProcessorWorkflow.csproj"
RUN dotnet build "OrderProcessorWorkflow.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "OrderProcessorWorkflow.csproj" -c Release -o /app/publish

FROM --platform=$TARGETPLATFORM mcr.microsoft.com/dotnet/aspnet:8.0 AS final
EXPOSE 80

ENV ASPNETCORE_URLS=http://+:80
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "OrderProcessorWorkflow.dll"]

# docker buildx build --platform inux/arm64 -t rochabr/order-processor-workflow:test .