<template>
  <div class="button" @click="$emit('click')">
    {{ title ? title : "" }}
    <font-awesome-icon v-if="icon" :icon="icon" slot="trigger" />
  </div>
</template>

<script>
export default {
  name: "Button",
  components: {},
  props: {
    title: String,
    icon: String,
  },
  data: function() {
    return {};
  },
};
</script>

<style scoped>
.button {
  display: inline-block;
  cursor: pointer;
  border-radius: 4px;
  background-color: silver;
  font-size: 13px;
  vertical-align: middle;
  padding: 8px;
  color: black;
  margin: inherit;
}

.button[size="small"] {
  padding: 4px;
}
.button[type="ok"] {
  background-color: #7957d5;
  color: white;
}

.button:hover {
  background-color: blue;
  color: white;
}
</style>
