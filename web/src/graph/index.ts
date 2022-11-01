import {ApolloClient, InMemoryCache} from "@apollo/client";

const Client = new ApolloClient({
	uri: "/api/v1/query",
	cache: new InMemoryCache()
});
export default Client;
