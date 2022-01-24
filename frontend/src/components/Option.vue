<template>
    <b-card v-bind:title="text">
    </b-card>
</template>


<script>
export default {
  name: 'Option',
  props: {
    text: String,
    selected: Boolean
  }
}
</script>