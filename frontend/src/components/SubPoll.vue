<template>
    <div class="pt-3">
        <b-card>
            <b-card-title>
                <b-row>
                    <b-col class="text-left">
                        {{ title }}
                    </b-col>
                    <b-col class="text-right">
                        <b-button v-bind:variant="!isOpen ? 'danger' : 'success'" v-on:click="toggle(id)">
                            {{ isOpen ? "Close" : "Open" }}
                        </b-button>
                    </b-col>
                </b-row>
            </b-card-title>
            <b-card-text>
            {{ description }}
            </b-card-text>
            <b-list-group>
                <b-list-group-item 
                button 
                v-for="option in options" 
                :key="option.id" 
                :disabled="!isOpen"
                :active="voted == option.id" 
                v-on:click="vote(id, option.id)"> {{ option.text }} </b-list-group-item>
                <b-list-group-item 
                button  
                :disabled="!isOpen"
                :active="voted == null"
                v-on:click="vote(id, null)"> Blank </b-list-group-item>
            </b-list-group>
        </b-card>
    </div>
</template>

<script>

export default {
  name: 'SubPoll',
  props: {
    id: String,
    title: String,
    description: String,
    options: Object,
    voted: String,
    isOpen: Boolean
  },
  methods: {
    vote: function(subpollId, optionId) {
        this.$store.commit("vote", {subpollId: subpollId, optionId: optionId});
    },
    toggle: function(subpollId) {
        this.$store.commit("toggleOpen", {subpollId: subpollId});
    }
  }
}
</script>