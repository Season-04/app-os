import { useQuery } from '@apollo/client'
import { ChevronRightIcon } from '@heroicons/react/20/solid'
import { graphql } from '../../gql'
import { UserRoleName } from './util'
import { Link, useNavigate, useParams } from 'react-router-dom'
import EditUserModal from './EditUserModal'
import { useState } from 'react'
import NewUserModal from './NewUserModal'

const getUsers = graphql(/* GraphQL */ `
  query GetUsers {
    users {
      id
      name
      emailAddress
      role
      lastSeenAt
      ...EditUserFragment
    }
  }
`)

export default function IndexPage() {
  const { userId } = useParams<{ userId: string }>()
  const navigate = useNavigate()

  const { loading, error, data, refetch } = useQuery(getUsers)
  console.log({ loading, error, data })

  const selectedUser =
    (userId && data?.users?.find((u) => u.id == userId)) || undefined

  const [showNewUserModal, setShowNewUserModal] = useState<boolean>(false)

  return (
    <div className="px-12 py-14">
      <div className="md:flex md:items-center md:justify-between pb-4">
        <div className="min-w-0 flex-1">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
            Users
          </h2>
        </div>
        <div className="mt-4 flex md:ml-4 md:mt-0">
          <button
            type="button"
            onClick={() => setShowNewUserModal(true)}
            className="ml-3 inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            New User
          </button>
        </div>
      </div>

      <ul role="list" className="divide-y divide-gray-100">
        {data?.users &&
          data.users.map((user) => (
            <li
              key={user.id}
              className="relative flex justify-between gap-x-6 py-5"
            >
              <div className="flex min-w-0 gap-x-4">
                {/* <img className="h-12 w-12 flex-none rounded-full bg-gray-50" src={person.imageUrl} alt="" /> */}
                <div className="min-w-0 flex-auto">
                  <p className="text-sm font-semibold leading-6 text-gray-900">
                    <Link to={`/settings/users/${user.id}`}>
                      <span className="absolute inset-x-0 -top-px bottom-0" />
                      {user.name}
                    </Link>
                  </p>
                  <p className="mt-1 flex text-xs leading-5 text-gray-500">
                    <a
                      href={`mailto:${user.emailAddress}`}
                      className="relative truncate hover:underline"
                    >
                      {user.emailAddress}
                    </a>
                  </p>
                </div>
              </div>
              <div className="flex shrink-0 items-center gap-x-4">
                <div className="hidden sm:flex sm:flex-col sm:items-end">
                  <p className="text-sm leading-6 text-gray-900">
                    {UserRoleName[user.role]}
                  </p>
                  {user.lastSeenAt && (
                    <p className="mt-1 text-xs leading-5 text-gray-500">
                      Last seen{' '}
                      <time dateTime={user.lastSeenAt}>{user.lastSeenAt}</time>
                    </p>
                  )}
                </div>
                <ChevronRightIcon
                  className="h-5 w-5 flex-none text-gray-400"
                  aria-hidden="true"
                />
              </div>
            </li>
          ))}
      </ul>

      {selectedUser && (
        <EditUserModal
          user={selectedUser}
          onClose={() => navigate('/settings/users')}
        />
      )}

      {showNewUserModal && !selectedUser && (
        <NewUserModal
          onClose={() => {
            refetch()
            setShowNewUserModal(false)
          }}
        />
      )}
    </div>
  )
}
