import { ApolloClient, HttpLink, InMemoryCache, from } from '@apollo/client';
import { onError } from '@apollo/client/link/error'

const errorLink = onError(({ graphQLErrors }) => {
  graphQLErrors?.forEach(e => {
    if (e.extensions['type'] === 'ACCESS_DENIED') {
      const loginUrl = new URL('/login', document.location.toString());
      loginUrl.searchParams.set('returnTo', document.location.pathname + document.location.search);
      document.location.assign(loginUrl);
    }
  })
})

const httpLink = new HttpLink({
  uri: '/api/core/graph',
})

export const client = new ApolloClient({
  link: from([errorLink, httpLink]),
  cache: new InMemoryCache(),
});
