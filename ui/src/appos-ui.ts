import { defineCustomElement, createApp } from 'vue'
import ApposChrome from './ApposChrome.ce.vue'
import './appos-ui.css'

// convert into custom element constructor
const ApposChromeElement = defineCustomElement(ApposChrome)

// register
customElements.define('appos-chrome', ApposChromeElement)
