import { ApolloClient, InMemoryCache } from '@apollo/client';

export const client = new ApolloClient({
  uri: '/api/core/graph',
  cache: new InMemoryCache(),
});
