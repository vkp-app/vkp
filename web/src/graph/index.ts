import {ApolloClient, InMemoryCache} from "@apollo/client";

const Client = new ApolloClient({
	uri: "/api/v1/query",
	cache: new InMemoryCache(),
	// disable caching so that
	// we always get up-to-date information
	// from the API
	defaultOptions: {
		query: {
			fetchPolicy: "no-cache",
			errorPolicy: "ignore"
		},
		watchQuery: {
			fetchPolicy: "no-cache",
			errorPolicy: "all"
		}
	}
});
export default Client;
