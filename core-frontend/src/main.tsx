import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import { client } from './apollo.ts'
import { ApolloProvider } from '@apollo/client'
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import LoginPage from './auth/LoginPage.tsx'
import UsersIndexPage from './users/IndexPage.tsx'


const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/settings/users",
    element: <UsersIndexPage />,
  },
  {
    path: "/settings/users/:userId",
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
