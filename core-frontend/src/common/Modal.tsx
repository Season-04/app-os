import { Dialog } from '@headlessui/react'

interface EditUserModalProps {
  title: string
  children: React.ReactNode
  buttons?: JSX.Element[]
  onClose: () => void
}

export default function Modal({
  title,
  children,
  buttons,
  onClose,
}: EditUserModalProps) {
  return (
    <Dialog open className="relative z-50" onClose={onClose}>
      <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />

      <div className="fixed inset-0 flex items-center justify-center p-4">
        <Dialog.Panel className="w-full max-w-xl overflow-hidden bg-white border border-zinc-300 rounded-xl">
          <Dialog.Title className="m-6 text-4xl font-bold">
            {title}
          </Dialog.Title>
          {children}
          {buttons && (
            <div className="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6">
              {buttons}
            </div>
          )}
        </Dialog.Panel>
      </div>
    </Dialog>
  )
}
