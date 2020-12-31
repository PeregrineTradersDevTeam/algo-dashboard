<template>
  <div>
    <responsive-modal
      name="window-launcher"
      :isActive.sync="isActive"
      :title="params.title"
      widthLimit="800"
    >
      <div v-if="params.isHolidayToday" class="holiday-hotify">
        Today is holiday! Please, pay attention to market holidays list
      </div>

      <div
        style="width:100%;height:40px;display:inline-block;margin:4px 2px 2px 2px;"
      >
        <b-field
          class="file is-primary"
          :class="{ 'has-name': !!file }"
          style="float:left"
        >
          <b-upload v-model="file" class="file-label" v-on:input="onFileChange">
            <span class="file-cta">
              <font-awesome-icon class="file" />
              <span class="file-label">Choose file</span>
            </span>
          </b-upload>
        </b-field>
        <b-field style="float:right; margin-right:5px">
          <b-button
            class="is-small is-success"
            @click="uploadFile()"
            :disabled="!file"
          >
            {{ params.buttonTitle }}
          </b-button>
        </b-field>
      </div>

      <div id="fileContent" class="f-content">
        <p v-if="fileContent != ''" v-html="fileContent"></p>
      </div>
    </responsive-modal>
  </div>
</template>

<script>
import { ToastProgrammatic as Toast } from "buefy";
import ResponsiveModal from "@/components/ResponsiveModal.vue";
import BField from "buefy/src/components/field/Field.vue";
import BUpload from "buefy/src/components/upload/Upload.vue";
import BButton from "buefy/src/components/button/Button.vue";
import BIcon from "buefy/src/components/icon/Icon.vue";
//, BUpload, BIcon } from "buefy";

export default {
  name: "LauncherModal",
  components: {
    ResponsiveModal: ResponsiveModal,
    "b-field": BField,
    "b-upload": BUpload,
    "b-button": BButton,
    "b-icon": BIcon,
  },

  data: function() {
    return {
      file: null,
      fileContent: "",
      isActive: false,
      params: {},
    };
  },
  methods: {
    showModal(title, buttonTitle, action, isHolidayToday) {
      this.isActive = true;
      this.params = { title, buttonTitle, action, isHolidayToday };
    },
    onFileChange(fileToLoad) {
      // console.log(fileToLoad);
      // console.log("onFileChange fired");

      var fileReader = new FileReader();

      let that = this;
      fileReader.onload = function(fileLoadedEvent) {
        that.fileContent = fileLoadedEvent.target.result.replace(
          /(\r\n|\n|\r)/g,
          "<br>"
        );
      };

      fileReader.readAsText(fileToLoad, "UTF-8");
    },

    uploadFile: function() {
      var data = new FormData();
      data.append("file", this.file);
      data.append("user", "rob");
      data.append("action", this.params.action);
      let that = this;

      fetch("/api/upload", {
        method: "POST",
        body: data,
      })
        .then(function(response) {
          if (!response.ok) {
            response.text().then(function(t) {
              Toast.open({
                message: t,
                type: "is-danger",
                position: "is-top",
              });
            });

            throw response;
          }
          return response;
        })
        .then((r) => r.json())
        .then(function(j) {
          let msg = "Configuration file ";
          if (j.received === true) {
            Toast.open({
              message: (msg += "received by trading server!"),
              type: "is-success",
              position: "is-top",
            });
            that.isActive = false;
            return;
          }

          Toast.open({
            message: (msg += " was not received by trading server!"),
            type: "is-danger",
            position: "is-top",
          });
        });
    },
  },
};
</script>

<style lang="scss" scoped>
@import "buefy/src/scss/buefy-build.scss";

.holiday-hotify {
  width: 100%;
  color: yellow;
  background-color: red;
  font-size: 11pt;
  font-weight: 800;
  margin-top: 0px;
  padding: 0;
  vertical-align: middle;
  padding: 0 10px 0 10px;
}
.file-name {
  border: 1px gray solid;
  color: rgba(255, 255, 255, 0.6);
}

.f-content {
  left: 1px;
  overflow-y: auto;
  display: block;
  height: 80%;
  margin: 4px 0 4px 2px;
  font-family: "Courier New", "Courier", monospace, monospace;
  color: rgba(255, 255, 255, 0.8);
  text-align: left;
  white-space: nowrap;
  overflow: auto;
}
</style>
