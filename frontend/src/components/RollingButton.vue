<template>
  <div class="floating-button" @click="next">
    <span class="label" v-if="mode == 'text'">{{
      title && title[selected] ? title[selected] : ""
    }}</span>
    <span class="icon" v-if="mode == 'icon'">
      <font-awesome-icon
        :icon="title && title[selected] ? title[selected] : ''"
        slot="trigger"
      />
    </span>
  </div>
</template>

<script>
export default {
  name: "RollingButton",
  props: {
    selected: {
      type: Number,
      default: 0,
    },
    title: Array,
    mode: {
      type: String,
      default: "text",
    },
    localStorageKey: {
      type: String,
      default: "",
    },
    KeyCode: String,
  },
  data: function() {
    return {};
  },
  methods: {
    next() {
      let idx = this.selected + 1;
      if (idx >= this.title.length) {
        idx = 0;
      }
      this.$emit("rolled", idx);
      if (this.localStorageKey) {
        localStorage.setItem(this.localStorageKey, idx.toString());
      }
    },
  },
  created() {
    if (!this.localStorageKey) {
      return;
    }
    const msm = localStorage.getItem(this.localStorageKey);
    if (!msm) {
      return;
    }
    let idx = parseInt(msm);
    if (idx < 0 || idx >= this.title.length) {
      return;
    }
    if (idx != this.selected) {
      this.$emit("rolled", idx);
    }
  },
  mounted() {
    if (!this.KeyCode) {
      return;
    }
    let that = this;
    window.addEventListener(
      "keypress",
      (e) => {
        switch (e.code) {
          case that.KeyCode:
            that.next();
            break;
          default:
            break;
        }
      },
      { passive: true }
    );
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.floating-button {
  right: 1px;
  height: 28px;
  width: 30px;
  cursor: pointer;
  border-radius: 4px;
  margin-left: 0px;
  background-color: silver;
  margin-bottom: 2px;
  margin-top: 2px;
  font-weight: bold;
  font-size: 12px;
  line-height: 12px;
  vertical-align: middle;
  color: rgba(0, 0, 0, 0.7);
}

.label {
  display: block;
  margin-top: 8px;
}

.icon {
  display: block;
  margin-top: 8px;
}

.floating-button:hover {
  background-color: blue;
  color: white;
}
</style>
