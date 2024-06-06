FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src
COPY . .
RUN dotnet restore "WorkerActorService.csproj"
RUN dotnet build "WorkerActorService.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "WorkerActorService.csproj" -c Release -o /app/publish

FROM --platform=$TARGETPLATFORM mcr.microsoft.com/dotnet/aspnet:8.0 AS final
EXPOSE 8080

ENV ASPNETCORE_URLS=http://+:8080
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "WorkerActorService.dll"]

