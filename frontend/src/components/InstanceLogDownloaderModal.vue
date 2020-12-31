<template>
  <responsive-modal
    name="window-instance-log-donwloade-modal"
    :isActive.sync="isActive"
    :title="id"
    widthLimit="80%"
  >
    <div style="display: flex; flex-direction: column; width:100%"></div>
  </responsive-modal>
</template>

<script>
import ResponsiveButton from "@/components/ResponsiveButton.vue";

export default {
  name: "InstanceLogDownloaderModal",
  components: {
    "responsive-modal": ResponsiveModal,
  },

  data: function() {
    return {
      isActive: false,
      id: "",
      activeTab: 0,
    };
  },
  methods: {
    showModal: function(prop, activeTab = 0) {
      this.id = prop.id;
      this.activeTab = activeTab;
      this.isActive = true;
    },

    hideModal: function() {
      this.isActive = false;
    },
    created: function() {
      this.isActive = false;
    },
  },
};
</script>

<style scoped></style>
