import { ChangeEventHandler } from "react"
import { graphql } from '../gql'
import { useMutation } from "@apollo/client"
import { useState } from "react"
import { UserRole } from "../gql/graphql"
import Modal from "../common/Modal"

export const CreateUserMutation = graphql(`
  mutation CreateUserMutation(
    $name: String!
    $emailAddress: String!
    $role: UserRole!
    $password: String!
  ) {
    createUser(input: {
      name: $name,
      emailAddress: $emailAddress,
      role: $role,
      password: $password
    }) {
      id
    }
  }
`)

interface NewUserModalProps {
  onClose: () => void;
}

export default function NewUserModal({ onClose }: NewUserModalProps) {
  const [createUser, { loading }] = useMutation(CreateUserMutation)

  const [name, setName] = useState<string>('')
  const [emailAddress, setEmailAddress] = useState<string>('')
  const [role, setRole] = useState<UserRole>(UserRole.User)
  const [password, setPassword] = useState<string>('')

  const onChangeRole: ChangeEventHandler<HTMLInputElement> = e => {
    setRole(e.target.value as UserRole)
  }

  const save = () => {
    createUser({ variables: { name, emailAddress, role, password } }).then(onClose)
  }

  const buttons = [
    <button
      type="button"
      className="inline-flex w-full justify-center rounded-md bg-red-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-red-500 sm:ml-3 sm:w-auto"
      disabled={loading}
      onClick={save}
    >
      Save
    </button>,
    <button
      type="button"
      className="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
      onClick={onClose}
    >
      Cancel
    </button>
  ]

  return <Modal
    title="New User"
    onClose={onClose}
    buttons={buttons}
  >
    <form className="space-y-6 m-6">
      <div className="flex flex-col gap-6 mt-8">
      <div className="flex flex-col">
          <label htmlFor="name" className="mb-2 font-medium leading-6 text-gray-900">
            Name
          </label>
          <input
            type="text"
            name="name"
            id="name"
            autoComplete='off'
            value={name}
            onChange={e => setName(e.target.value)}
            className="w-full rounded-md border-0 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-sky-600 sm:py-1.5 sm:text-sm sm:leading-6"
          />
        </div>

        <div className="flex flex-col">
          <label htmlFor="emailAddress" className="mb-2 font-medium leading-6 text-gray-900">
            Email Address
          </label>
          <input
            type="email"
            name="emailAddress"
            id="emailAddress"
            autoComplete='off'
            value={emailAddress}
            onChange={e => setEmailAddress(e.target.value)}
            className="w-full rounded-md border-0 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-sky-600 sm:py-1.5 sm:text-sm sm:leading-6"
          />
        </div>

        <div className="flex flex-col">
          <legend className="mb-2 font-medium leading-6 text-gray-900">Role</legend>
          <div className="ml-2 flex flex-col">
            <label>
              <input type="radio" name="role" value="ADMIN" checked={role == 'ADMIN'} onChange={onChangeRole} /> Admin
            </label>
            <label>
              <input type="radio" name="role" value="USER" checked={role == 'USER'} onChange={onChangeRole} /> User
            </label>
          </div>
        </div>


        <div className="flex flex-col">
          <label htmlFor="password" className="mb-2 font-medium leading-6 text-gray-900">
            Password
          </label>
          <input
            type="password"
            name="password"
            id="password"
            autoComplete='off'
            value={password}
            onChange={e => setPassword(e.target.value)}
            className="w-full rounded-md border-0 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-sky-600 sm:py-1.5 sm:text-sm sm:leading-6"
          />
        </div>
      </div>
    </form>
  </Modal>
}
