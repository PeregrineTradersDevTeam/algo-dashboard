<template>
  <div class="tab-content-console">
    <div class="subheader">
      <button
        class="button is-info"
        :disabled="!isShowNewVisible"
        v-on:click="getNewLogItems()"
      >
        Show new
      </button>
    </div>
    <div class="aa">
      <Table
        :rows="items"
        :head="['Date', 'Time', 'Level', 'Message']"
        :colstyle="{
          0: 'width: 75px',
          1: 'width: 80px',
          2: 'width: 40px',
        }"
        :cellstyle="{
          2: (val) => {
            return (
              'background-color: ' +
              config.severityColor[val] +
              ';color: black;'
            );
          },
        }"
        :replacer="{
          2: (val) => {
            return config.severityLabel[val];
          },
        }"
      >
        <template v-slot:footer>
          <button
            v-show="isShowMoreVisible"
            v-on:click="getNextLogItems(config.logPerPage)"
          >
            Show previous
          </button>
          <button v-show="isShowMoreVisible" v-on:click="getAllLogItems()">
            Show All
          </button>
        </template>
      </Table>
    </div>
  </div>
</template>

<script>
function formatD(date) {
  var d = new Date(date),
    month = "" + (d.getMonth() + 1),
    day = "" + d.getDate(),
    year = d.getFullYear();

  if (month < 10) month = "0" + month;
  if (day < 10) day = "0" + day;

  return [year, month, day].join("-");
}

function formatT(date) {
  var d = new Date(date),
    h = d.getHours(),
    m = d.getMinutes(),
    s = d.getSeconds(),
    ms = d.getMilliseconds();

  if (h < 10) h = "0" + h;
  if (m < 10) m = "0" + m;
  if (s < 10) s = "0" + s;
  if (ms < 10) {
    ms = "00" + ms;
  }
  if (ms >= 10 && ms < 100) {
    ms = "0" + ms;
  }

  return [h, m, s].join(":") + "." + ms;
}

async function fetchLogItems(id, from, to, filter) {
  let url = "/api/log/" + id + "/feed?from=" + from + "&to=" + to;
  if (filter) {
    url += "&filter=" + filter;
  }

  const resp = await fetch(url);
  var data = await resp.json();

  let res = [];
  for (var i in data) {
    //    console.log("tmp.i=", i);
    // console.log("r: ", res[i]);
    res.push({
      atd: formatD(data[i].at),
      att: formatT(data[i].at),
      severity: data[i].severity,
      msg: data[i].msg,
    });
    //res.at = formatDate(res[i].at);
    //console.log("res.lentht in for: ", res.length);
  }
  // console.log("res: ", res);

  return res;
}
import Table from "@/components/Table.vue";

export default {
  name: "InstanceModalLog",
  components: {
    Table,
  },
  props: {
    code: String,
    filter: {
      type: String,
      default: "",
    },
  },
  data: function() {
    return {
      id: "",
      items: [],
      isShowMoreVisible: false,
      isShowNewVisible: false,
      totalLen: 0,
      newTotalLen: 0,
      polling: null,
      config: {
        logPerPage: 20,
        severityColor: {
          0: "red",
          1: "orange",
          2: "yellow",
          3: "lightgreen",
        },
        severityLabel: {
          0: "FATL",
          1: "ERRO",
          2: "WARN",
          3: "INFO",
        },
      },
    };
  },
  computed: {
    sevcolor() {
      return "red";
    },
  },
  methods: {
    getLogItems(id) {
      let that = this;

      this.id = id;
      this.isShowMoreVisible = false;

      this.getCount(id).then((cnt) => {
        that.totalLen = cnt;

        fetchLogItems(id, 0, this.config.logPerPage - 1, this.filter).then(
          (arr) => {
            that.items.push(...arr);
            that.isShowMoreVisible = that.items.length < that.totalLen;
          }
        );

        that.polling = setInterval(function() {
          that.getCount(id).then((cnt) => {
            that.newTotalLen = cnt;
            if (that.newTotalLen > that.totalLen) {
              that.isShowNewVisible = true;
            }
          });

          /*console.log(
            "newTotalLen",
            that.newTotalLen,
            "totalLen",
            that.totalLen
          );*/
        }, 5000);
      });

      //console.log("Entering getLogItems(id), totalLen=", id, that.totalLen);
    },

    getNextLogItems(cnt) {
      let fincnt = -1;
      if (cnt != -1) {
        fincnt = this.items.length + cnt - 1;
      }
      fetchLogItems(this.id, this.items.length, fincnt).then((v) =>
        this.items.push(...v)
      );
    },

    getNewLogItems() {
      let to = this.newTotalLen - this.totalLen;
      if (to == 0) {
        return;
      }
      fetchLogItems(this.id, 0, to).then((arr) => {
        this.items.unshift(...arr);
        this.totalLen = this.newTotalLen;
        this.isShowNewVisible = false;
      });
    },

    getAllLogItems() {
      let that = this;
      fetchLogItems(this.id, 0, -1).then((arr) => {
        that.items = [];
        that.items.push(...arr);
        that.totalLen = arr.length;
        that.isShowMoreVisible = false;
      });
    },

    // getCount return total length of the list L:XXXX where XXXX is instance code.
    getCount(id) {
      return fetch("/api/log/" + id + "/count")
        .then((r) => r.json())
        .then((j) => {
          return j.count;
        });
    },
  },
  mounted() {
    //console.log("before updated instance modal info");
    this.getLogItems(this.$props.code);
  },
  beforeDestroy() {
    clearInterval(this.polling);
  },
};
</script>

<style scoped>
.subheader {
  display: flex;
  position: sticky;
  top: 0px;
  background: #1e1e1e;
  height: 25px;
  justify-content: space-between;
  font-size: 120%;
}

.tab-content-console {
  display: block;
  width: 100%;
  height: calc(100% - 2px);
  background-color: #1e1e1e;
  overflow-y: hidden;
}

.footer {
  display: flex;
  position: sticky;
  bottom: 0px;
  background: #1e1e1e;
  height: 25px;
  justify-content: space-between;
  font-size: 120%;
}
.aa {
  display: block;
  height: calc(100% - 40px);
  top: 0px;
  margin: 0;
  background-color: #1e1e1e;
}
</style>
