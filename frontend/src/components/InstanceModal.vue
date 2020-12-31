<template>
  <responsive-modal
    name="window-instance-modal"
    :isActive.sync="isActive"
    :title="id"
    widthLimit="80%"
  >
    <div style="display: flex; flex-direction: column; width:100%">
      <Tabs :title="['console', 'inputs', 'orders']">
        <template v-slot:console>
          <InstanceModalLog :code="id" />
        </template>
        <template v-slot:inputs>
          <InstanceModalInfo :code="id" />
        </template>
        <template v-slot:orders>
          <InstanceModalOrders :code="id" />
        </template>
      </Tabs>
    </div>
  </responsive-modal>
</template>

<script>
import InstanceModalInfo from "@/components/InstanceModalInfo.vue";
import InstanceModalLog from "@/components/InstanceModalLog.vue";
import ResponsiveModal from "@/components/ResponsiveModal.vue";
import ResponsiveButton from "@/components/ResponsiveButton.vue";
import InstanceModalOrders from "@/components/InstanceModalOrders.vue";

//import { Tabs, Tab } from "vue-slim-tabs";
//import "buefy/dist/buefy.css";
import Tabs from "@/components/Tabs.vue";

export default {
  name: "InstanceModal",
  components: {
    InstanceModalInfo: InstanceModalInfo,
    InstanceModalLog: InstanceModalLog,
    "responsive-modal": ResponsiveModal,
    InstanceModalOrders: InstanceModalOrders,
    Tabs: Tabs,
    //tab: Tab,
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
