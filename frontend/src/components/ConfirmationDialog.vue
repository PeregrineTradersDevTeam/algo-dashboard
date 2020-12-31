<template>
  <responsive-modal
    name="window-config"
    :isActive.sync="pressed"
    title="Settings"
    widthLimit="60%"
  >
    <!-- <template v-slot:default>
      <b-field horizontal label="Update,sec" style="color:white">
        <b-numberinput
          v-model="config.refreshInterval"
          min="1"
          size="is-small"
          controls-position="compact"
          icon-pack="fa"
        ></b-numberinput>
      </b-field>
      <b-checkbox v-model="config.isRiskVisible"
        >Enable Risk Information</b-checkbox
      >
    </template>
    <template v-slot:footer>
      <button
        class="button is-small "
        @click="
          {
            $emit('save', config);
            pressed = false;
          }
        "
      >
        Save
      </button>
      <button class="button is-small" @click="pressed = false">
        Cancel
      </button>
    </template> -->
  </responsive-modal>
</template>

<script>
import ResponsiveModal from "@/components/ResponsiveModal.vue";

export default {
  name: "ConfirmationDialog",
  components: {
    "responsive-modal": ResponsiveModal,
  },
  props: {
    buttons: Array,
    confirmIndex: Number,
  },
  data: function() {
    return {
      pressed: false,
      config: {
        refreshInterval: 5,
        isRiskVisible: false,
      },
    };
  },
  methods: {
    show: function(cfg) {
      this.config.isRiskVisible = cfg.isRiskVisible;
      this.config.refreshInterval = cfg.refreshInterval;
      this.pressed = true;
    },
    save: function() {
      this.pressed = false;
      return {
        isRiskVisible: this.config.isRiskVisible,
        refreshInterval: this.config.refreshInterval,
      };
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.floating-button {
  height: 24px;
  width: 28px;
  padding: 1px 1px 0px 1px;
  font-weight: bold;
  font-size: 14px;
  border: 1px solid silver;
  cursor: pointer;
  border-radius: 4px;
  margin-left: 0px;
  margin-top: 2px;
  background-color: silver;
}

.floating-button:hover {
  background-color: blue;
  color: white;
}
</style>
