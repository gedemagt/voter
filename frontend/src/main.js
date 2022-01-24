import '@babel/polyfill'
import 'mutationobserver-shim'
import Vue from 'vue'
import './plugins/bootstrap-vue'
import App from './App.vue'
import Vuex from 'vuex'
import { BootstrapVue, BootstrapVueIcons } from 'bootstrap-vue'

Vue.use(Vuex)
Vue.config.productionTip = false
Vue.use(BootstrapVue)
Vue.use(BootstrapVueIcons)

// Create a new store instance.
const store = new Vuex.Store({
  state: {
    subpolls: {
      "1": {
        id: "1",
        title: "SubPoll 1",
        description: "Description 1",
        voted: null,
        isOpen: true,
        options: {
         "1": {
           "id": "1",
           "text": "Option1"
         },
         "2": {
            "id": "2",
            "text": "Option2"
         },
         "3": {
            "id": "3",
            "text": "Option3"
         },
         "4": {
            "id": "4",
            "text": "Option4"
         }
        }
      },
      "2": {
        id: "2",
        title: "SubPoll 2",
        description: "Description 2",
        voted: null,
        isOpen: false,
        options: {
         "1": {
           "id": "1",
           "text": "Option1"
         },
         "2": {
            "id": "2",
            "text": "Option2"
         },
         "3": {
            "id": "3",
            "text": "Option3"
         },
         "4": {
            "id": "4",
            "text": "Option4"
         }
        }
      }
    }
  },
  mutations: {
    vote (state, payload) {
      state.subpolls[payload.subpollId].voted = payload.optionId;
    },
    toggleOpen (state, payload) {
      state.subpolls[payload.subpollId].isOpen = !state.subpolls[payload.subpollId].isOpen;
    }
  }
})


var ws = new WebSocket("ws://localhost:8090/echo");
ws.onopen = function() {
                  
  // Web Socket is connected, send data using send()
  ws.send("Message to send");
  alert("Message is sent...");
};

ws.onmessage = function (evt) { 
  var received_msg = evt.data;
  alert("Message is received..." + received_msg);
};

ws.onclose = function() { 
  
  // websocket is closed.
  alert("Connection is closed..."); 
};


new Vue({
  render: h => h(App),
  store: store,
}).$mount('#app')
