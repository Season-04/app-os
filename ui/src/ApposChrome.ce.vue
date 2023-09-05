<template>
  <div class="h-full">
    <div class="fixed inset-y-0 z-50 flex w-72 flex-col">
      <div class="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-900 px-6">
        <div class="flex h-16 shrink-0 items-center">
          <img
            class="h-8 w-auto"
            src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=500"
            alt="Your Company"
          />
        </div>
        <nav class="flex flex-1 flex-col">
          <ul role="list" class="flex flex-1 flex-col gap-y-7">
            <li>
              <ul role="list" class="-mx-2 space-y-1">
                <li v-for="item in navigation" :key="item.href">
                  <a
                    :href="item.href"
                    :class="[
                      item.current
                        ? 'bg-gray-800 text-white'
                        : 'text-gray-400 hover:text-white hover:bg-gray-800',
                      item.open && 'text-white',
                      'group flex items-center gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold'
                    ]"
                  >
                    <component :is="item.icon" class="h-6 w-6 shrink-0" aria-hidden="true" />
                    {{ item.name }}
                    <ChevronRightIcon
                      v-if="item.children"
                      :class="[
                        item.open ? 'rotate-90 text-gray-400' : 'text-gray-500',
                        'ml-auto h-5 w-5 shrink-0'
                      ]"
                      aria-hidden="true"
                    />
                  </a>
                  <ul v-if="item.open && item.children" role="list" class="space-y-1">
                    <li v-for="child in item.children" :key="child.href">
                      <a
                        :href="child.href"
                        :class="[
                          child.current
                            ? 'bg-gray-800 text-white'
                            : 'text-gray-400 hover:text-white hover:bg-gray-800',
                          'group flex gap-x-3 rounded-md p-2 pl-11 text-sm leading-6 font-semibold'
                        ]"
                      >
                        {{ child.name }}
                      </a>
                    </li>
                  </ul>
                </li>
              </ul>
            </li>
            <li class="-mx-6 mt-auto">
              <a
                href="#"
                class="flex items-center gap-x-4 px-6 py-3 text-sm font-semibold leading-6 text-white hover:bg-gray-800"
              >
                <img
                  class="h-8 w-8 rounded-full bg-gray-800"
                  src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                  alt=""
                />
                <span class="sr-only">Your profile</span>
                <span aria-hidden="true">Tom Cook</span>
              </a>
            </li>
          </ul>
        </nav>
      </div>
    </div>

    <main class="pl-72 h-full">
      <slot />
    </main>
  </div>
</template>

<style scoped>
@tailwind base;
@tailwind components;
@tailwind utilities;
</style>

<script setup lang="ts">
import {
  ChevronRightIcon,
  HomeIcon,
  UsersIcon,
  ClockIcon,
  AdjustmentsHorizontalIcon
} from '@heroicons/vue/24/outline'

const pathNameWithoutTrailingSlash = document.location.pathname.replace(/\/$/, '')

const isCurrentRoute = (path: string, exact = false) => {
  if (document.location.pathname === path) {
    return true
  }

  if (pathNameWithoutTrailingSlash === path) {
    return true
  }

  return !exact && pathNameWithoutTrailingSlash.startsWith(path)
}

const navigation = [
  {
    name: 'Dashboard',
    href: '/',
    icon: HomeIcon,
    current: document.location.pathname === '/'
  },
  {
    name: 'Clock',
    href: '/clock',
    icon: ClockIcon,
    current: isCurrentRoute('/clock')
  },
  {
    name: 'Settings',
    href: '/settings',
    icon: AdjustmentsHorizontalIcon,
    current: isCurrentRoute('/settings', true),
    open: isCurrentRoute('/settings', false),
    children: [
      {
        name: 'Users',
        href: '/settings/users',
        current: isCurrentRoute('/settings/users')
      },
      {
        name: 'Applications',
        href: '/settings/applications',
        current: isCurrentRoute('/settings/applications')
      }
    ]
  }
]
</script>
