import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import { client } from './apollo.ts'
import { ApolloProvider } from '@apollo/client'
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import UsersIndexPage from './users/IndexPage.tsx'

const router = createBrowserRouter([
  {
    path: "/users",
    element: <UsersIndexPage />,
  },
  {
    path: "/users/:userId",
    element: <UsersIndexPage />
  }
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ApolloProvider client={client}>
      <RouterProvider router={router} />
    </ApolloProvider>
  </React.StrictMode>
)
