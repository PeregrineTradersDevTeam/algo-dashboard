<template>
  <section>
    <b-modal :active.sync="isActive">
      <header class="modal-card-head">
        <p class="modal-card-title">Platform</p>
        <button class="delete" aria-label="close" @click="hideModal()"></button>
      </header>
      <section class="modal-card-body">
        <b-tabs v-model="activeTab" @input="switchTab()" destroy-on-hide>
          <b-tab-item label="1 - Console" style="height:400px; overflow-y:auto">
            <InstanceModalLog :code="id" />
          </b-tab-item>
        </b-tabs>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-primary" @click="hideModal()">
          Hide window
        </button>
      </footer>
    </b-modal>
  </section>
</template>

<script>
//import InstanceModalInfo from "@/components/InstanceModalInfo.vue";
import InstanceModalLog from "@/components/InstanceModalLog.vue";
//import "buefy/dist/buefy.css";

export default {
  name: "PlatformModal",
  components: {
    //InstanceModalInfo: InstanceModalInfo,
    InstanceModalLog: InstanceModalLog,
  },

  data: function() {
    return {
      isActive: false,
      id: "",
      activeTab: 0,
    };
  },
  methods: {
    showModal: function(prop) {
      this.id = prop.id;
      this.isActive = true;
      //  console.log("I am here", prop, this);
    },

    switchTab: function() {},

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
