<template>
  <div
    id="status-labels"
    v-bind:style="{
      display: 'flex',
      flexDirection: location === 'header' ? 'row' : 'column',
    }"
  >
    <span
      v-for="(cnt, status) in instanceStatusCounters"
      v-bind:key="status"
      class="header-counter"
      v-bind:style="{ 'background-color': config.instanceStatusColor[status] }"
    >
      {{ cnt == 0 ? "&centerdot;" : cnt }}
    </span>

    <!-- <div class="buttons">
      <b-button
        class="is-small is-success short-button-text"
        v-on:click="launch"
        >L</b-button
      >
      <b-button
        class="is-small is-danger short-button-text"
        v-on:click="confirm"
        >C</b-button
      >
      <b-button class="is-small is-success full-button-text" v-on:click="launch"
        >Launcher</b-button
      >
      <b-button class="is-small full-button-text" v-on:click="update"
        >Update</b-button
      >
      <b-button class="is-small is-danger full-button-text" v-on:click="confirm"
        >Close ALL</b-button
      >
    </div>
    <LauncherModal ref="launcher" /> -->
  </div>
</template>

<script>
import { ToastProgrammatic as Toast } from "buefy";
import LauncherModal from "@/components/LauncherModal.vue";

export default {
  name: "InstanceStatusCounter",
  components: {
    LauncherModal,
  },
  props: {
    location: String,
  },
  data: function() {
    return {
      config: {
        instanceStatusColor: {},
      },
      isCalendarActive: false,
    };
  },
  computed: {
    instanceStatusCounters: function() {
      return this.$store.getters.instanceStatusCounter;
    },
  },

  methods: {
    statusColor: function(status) {
      return this.config.instanceStatusColor.get(status);
    },
    launch: function() {
      this.$refs.launcher.showModal(
        "Algo Launcher",
        "Start ALL",
        "launch",
        this.$store.getters.holidayToday
      );
    },

    update: function() {
      this.$refs.launcher.showModal(
        "Algo Updater",
        "Send Update",
        "update",
        this.$store.getters.holidayToday
      );
    },

    confirm: function() {
      this.$buefy.dialog.confirm({
        message: "Close all instances?",
        focusOn: "confirm",
        animation: "",
        type: "is-danger",
        onConfirm: () => this.closeAll(),
      });
    },

    closeAll: function() {
      fetch("/api/close-all", { method: "POST" })
        .then((r) => r.json())
        .then(function(j) {
          if (j.received === true) {
            Toast.open({
              message: "<strong>CLOSE ALL</strong> received by server!",
              type: "is-success",
              position: "is-top-right",
            });
          } else {
            Toast.open({
              message:
                "<strong>CLOSE ALL</strong>  has not received by server!",
              type: "is-danger",
              position: "is-top-right",
            });
          }
        });
    },
  },

  created: function() {
    this.config.instanceStatusColor = this.$store.getters.instanceStatusColor;
  },
  mounted: function() {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.horizontal {
}

/*
@media only screen and (max-width: 1040px) {
  .horizontal {
    display: inline;
  }
  .vertical {
    display: none;
  }
}
@media only screen and (max-width: 776px) {
  .horizontal {
    display: inline;
  }
  .vertical {
    display: none;
  }
}

@media only screen and (max-width: 563px) {
  .horizontal {
    display: none;
  }
  .vertical {
    display: block;
    position: sticky;
    top: 27px;
    right: 1px;
  }
} */

.header-counter {
  display: block;
  height: 22px;
  width: 27px;
  padding: 1px 1px 1px 1px;
  margin: 0px 1px 1px 1px;
  font-weight: bold;
  font-size: 14px;
  line-height: 22px;
  color: black;
}
/*
@media only screen and (max-width: 563px) {
  .header-counter {
    margin-top: -1px;
    height: 24px;
    width: 28px;
    padding: 2px 1px 2px 1px;
    font-weight: bold;
    font-size: 14px;
    border: 1px;
    border-style: solid;
    border-color: black;
    margin-left: 1px;
  }
  .button-calendar {
    margin-top: 5px;
    height: 24px;
    width: 28px;
    padding: 2px 1px 2px 1px;
    font-weight: bold;
    font-size: 14px;
    border: 1px solid silver;
    margin-left: 1px;
  }
}

@media only screen and (max-width: 894px) {
  .short-button-text {
    display: inline-block;
  }
  .full-button-text {
    display: none;
  }
  .tight {
    padding-left: 6px;
    padding-right: 6px;
  }
} */

.portfolio-pnl {
  font-size: 12px;
  font-weight: bold;
}

.buttons {
  flex-direction: row;
  flex-wrap: nowrap;
  flex: 1;
  align-content: flex-start;
}
.separator {
  width: 20px;
}

.short-button-text {
  display: none;
}

.vertical {
  display: block;
  position: sticky;
  top: 27px;
  right: 1px;
}
</style>
