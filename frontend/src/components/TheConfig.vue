<template>
  <responsive-modal
    name="window-config"
    :isActive.sync="pressed"
    title="Settings"
    widthLimit="60%"
    :resizeable="false"
  >
    <template v-slot:default>
      <div class="field">
        <label for="mrf">Matrix, refresh interval</label>
        <input
          id="mrf"
          type="number"
          v-model="config.refreshInterval"
          ref="mrf"
          min="1"
          max="60"
          step="1"
        />
        <label for="mrf">(seconds)</label>
      </div>
      <div class="field">
        <label for="eri">Enable risk information</label>
        <input id="eri" type="checkbox" v-model="config.isRiskVisible" />
      </div>
      <div class="button-group">
        <Button
          type="ok"
          @click="
            {
              $emit('save', save());
              pressed = false;
            }
          "
          title="Save"
        />
        <Button
          class="button is-small"
          @click="pressed = false"
          title="Cancel"
        />
      </div>
    </template>
    <template v-slot:footer> </template>
  </responsive-modal>
</template>

<script>
import ResponsiveModal from "@/components/ResponsiveModal.vue";
import Button from "@/components/Button.vue";

export default {
  name: "TheConfig",
  components: {
    "responsive-modal": ResponsiveModal,
    Button: Button,
  },
  props: {
    cfg: Object,
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
    showModal(cfg) {
      this.config.isRiskVisible = cfg.isRiskVisible;
      this.config.refreshInterval = cfg.refreshInterval;
      this.pressed = true;
      this.$nextTick(() => this.$refs.mrf.focus());
    },
    save() {
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
.field {
  display: block;
  text-align: left;
  padding: 10px;
  font-size: 13px;
}
.field label {
  display: inline-block;
  margin-right: 10px;
  white-space: nowrap;
  vertical-align: middle;
}

.field input[type="number"] {
  font-size: 14px;
  padding: 4px;
}
.field input[type="checkbox"] {
  display: inline-block;
  transform: scale(1.3);
  vertical-align: middle;
}
.button-group {
  margin-right: 8px;
}
</style>
