import { ChangeEventHandler } from "react"
import { graphql } from '../gql'
import { ResultOf } from "@graphql-typed-document-node/core"
import { useMutation } from "@apollo/client"
import { useState } from "react"
import { UserRole } from "../gql/graphql"
import Modal from "../common/Modal"

export const EditUserFragment = graphql(/* GraphQL */ `
  fragment EditUserFragment on User {
    id
    name
    emailAddress
    role
  }
`)

export const UpdateUserMutation = graphql(`
  mutation UpdateUserMutation(
    $id: ID!,
    $name: String!
    $role: UserRole!
  ) {
    updateUser(input: {id: $id, name: $name, role: $role}) {
      ...EditUserFragment
    }
  }
`)

type User = ResultOf<typeof EditUserFragment>

interface EditUserModalProps {
  user: User;
  onClose: () => void;
}

export default function EditUserModal({ user, onClose }: EditUserModalProps) {
  console.log('EditUserModal', { user })
  const [updateUser, { loading }] = useMutation(UpdateUserMutation)

  const [name, setName] = useState<string>(user?.name || '')
  const [role, setRole] = useState<UserRole>(user?.role || UserRole.User)

  const onChangeRole: ChangeEventHandler<HTMLInputElement> = e => {
    setRole(e.target.value as UserRole)
  }

  const update = () => {
    updateUser({ variables: { id: user!.id, name, role } }).then(onClose)
  }

  const buttons = [
    <button
      type="button"
      className="inline-flex w-full justify-center rounded-md bg-red-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-red-500 sm:ml-3 sm:w-auto"
      disabled={loading}
      onClick={update}
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
    title={`Edit ${user.name}`}
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
      </div>
    </form>
  </Modal>
}
