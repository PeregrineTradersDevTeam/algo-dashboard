<template>
  <div>
    <responsive-modal
      name="window-launcher"
      :isActive.sync="isActive"
      :title="params.title"
      widthLimit="800"
      @close="clean"
    >
      <div v-if="params.isHolidayToday" class="holiday-banner">
        Attention: Today is holiday!
      </div>

      <div class="toolbar">
        <div style="display: flex; ">
          <ResponsiveButton
            title="Choose file"
            bgColor="#7957d5"
            color="white"
            icon="upload"
            @click="$refs.upload.click()"
          />
          <input type="file" ref="upload" hidden @change="onFileChange" />
          <span class="file-name">{{ filename }}</span>
        </div>
        <div style="display: flex">
          <ResponsiveButton
            :title="params.buttonTitle"
            bgColor="#48c774"
            color="white"
            @click="uploadFile()"
            :disabled="!fileContent"
          />
        </div>
      </div>

      <div
        id="fileContent"
        class="f-content"
        v-bind:style="{
          height: params.isHolidayToday
            ? 'calc(100% - 55px)'
            : 'calc(100% - 35px)',
        }"
      >
        <p v-if="fileContent != ''" v-html="fileContent"></p>
      </div>
    </responsive-modal>
  </div>
</template>

<script>
import ResponsiveModal from "@/components/ResponsiveModal.vue";
import ResponsiveButton from "@/components/ResponsiveButton.vue";

export default {
  name: "LauncherModal",
  components: {
    ResponsiveModal,
    ResponsiveButton,
  },

  data: function() {
    return {
      filename: "No file choosen",
      file: null,
      fileContent: "",
      params: {},
      isActive: false,
    };
  },
  methods: {
    clean() {
      this.filename = "No file choosen";
      this.fileContent = "";
      this.file = null;
    },
    showModal(title, buttonTitle, action, isHolidayToday) {
      this.isActive = true;
      this.params = { title, buttonTitle, action, isHolidayToday };
    },
    onFileChange(event) {
      const files = event.target.files || event.dataTransfer.files;
      if (!files || files.length == 0) {
        return;
      }

      this.file = files[0];
      this.filename = this.file.name;

      var fileReader = new FileReader();

      let that = this;
      fileReader.onload = function(fileLoadedEvent) {
        that.fileContent = fileLoadedEvent.target.result.replace(
          /(\r\n|\n|\r)/g,
          "<br>"
        );
      };

      fileReader.readAsText(this.file, "UTF-8");
    },

    uploadFile: function() {
      let data = new FormData();
      data.append("file", this.file, this.file.name);
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
              that.$toast.error(t);
            });
            throw response;
          }
          return response;
        })
        .then((r) => r.json())
        .then(function(j) {
          let msg = "Configuration file ";
          if (j.received === true) {
            that.$toast.success(msg + " received by trading server!");
            that.isActive = false;
            that.clean();
            return;
          }
          that.$toast.error(msg + " was not received by trading server!");
        });
    },
  },
};
</script>

<style scoped>
.holiday-banner {
  color: yellow;
  background-color: red;
  font-size: 11pt;
  line-height: 20px;
  height: 20px;
  font-weight: 800;
  margin-top: 0px;
  padding: 0;
  vertical-align: middle;
  padding: 0 10px 0 10px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  padding: 4px;
  font-size: 140%;
  background: rgba(0, 0, 0, 0.3);
}

.file-name {
  display: flex;
  align-items: center;
  margin-left: 10px;
}

.f-content {
  display: block;
  margin: 4px 0 4px 2px;
  font-family: "Courier New", "Courier", monospace, monospace;
  color: rgba(255, 255, 255, 0.8);
  text-align: left;
  font-size: 120%;
  white-space: nowrap;
  overflow: auto;
}
</style>
