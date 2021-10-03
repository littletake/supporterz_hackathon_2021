<template>
  <div class="Input">
    <div class="header">
      <div class="title">
        <div class="headerTitle" v-on:click="onClickTitle">フォト旅</div>
      </div>
      <div class="menu">
        <div class="toHome" v-on:click="onClickTitle">Home</div>
        <div class="userSetting" v-on:click="onClickOpenPopup">
          {{ username_upper }}
        </div>
      </div>
      <Popup
        :isopen="popup"
        :user="username"
        :tripnum="tripIdNum"
        @close="onClickClosePopup"
      ></Popup>
    </div>
    <div class="input_body">
      <div class="input_area">
        <!-- <input type="file" name="example" ref="preview" accept="image/*" multiple required> -->
        <p id="error" v-show="error">{{ error }}</p>
        <label>
          <!-- <img :src="avatar" alt="Avatar" class="image"> -->
          <div
            class="drop_area"
            @dragenter="dragEnter"
            @dragleave="dragLeave"
            @dragover.prevent
            @drop.prevent="dropFile()"
            :class="{ enter: isEnter }"
            @change="onImageChange"
          >
            Drop files here or click to upload
            <!-- <div> -->

            <input
              type="file"
              accept="image/jpeg"
              @change="onImageChange"
              multiple
            />
          </div>
        </label>
        <div id="howToUse">JPEG, PNG format can use this application</div>

        <br />
        <div>
          <ul class="flex">
            <li
              class="flex-col"
              v-for="(file, index) in files"
              :key="index"
              @click="deleteFile(index)"
            >
              <!-- {{ index }} -->
              <div style="position: relative">
                <span class="delete-mark">×</span>
                <img class="file_icon" src="../assets/icon.png" />
                <!-- <img class="file_icon" :src="images[index]"> -->
              </div>
              <span>{{ file.name }}</span>
            </li>
          </ul>
        </div>

        <br />
        <div v-show="files.length">
          <button class="button" v-on:click="upload" v-bind:disabled="isPushed">
            {{ button_text }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import loadImage from "blueimp-load-image";
import Popup from "@/components/UserPopup.vue";

export default {
  components: { Popup },
  name: "Input",
  data() {
    return {
      error: "",
      isEnter: false,
      files: [],
      images: [],
      test: "",
      isPushed: false,
      button_text: "送信",
      username: "undefined user",
      tripIdNum: "",
      popup: false,
      username_upper: "",
    };
  },
  mounted: function () {
    this.tripIdNum = this.$store.getters["trip/getNumTripID"];
    this.username = this.$store.getters["user/getUserName"];
    this.username_upper = this.username.slice(0, 1).toUpperCase();
  },

  methods: {
    // setError (error, text) {
    //   this.error = (error.response && error.response.data && error.response.data.error) || text
    // },
    changeButton() {
      this.isPushed = true;
      this.button_text = "送信中";
    },
    dragEnter() {
      this.isEnter = true;
    },
    dragLeave() {
      this.isEnter = false;
    },
    deleteFile(index) {
      this.files.splice(index, 1);
      this.images.splice(index, 1);
    },

    checkTypeAndGPS(files) {
      return new Promise((resolve) => {
        files.forEach((file) => {
          // 拡張子のチェック
          const type = file.type;
          if (type != "image/jpeg") {
            alert(
              "エラー: " + file.name + "\n拡張子がjpegの画像を選択してください"
            );
            return;
          }
          // GPSのチェック
          loadImage.parseMetaData(file, (data) => {
            if (data.exif && data.exif.get("GPSInfo")) {
              // console.log("fileWithGPS ", fileWithGPS);
              // console.log("3: GPS取得できてる");
              this.processFile(file);
            } else {
              alert(
                "エラー: " +
                  file.name +
                  "\n位置情報を含む画像を選択してください"
              );
            }
          });
        });
        resolve();
      });
    },

    processFile(file) {
      // console.log("4: imagesを抜き出す処理")
      // console.log("5: ",file)
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        const im = reader.result;
        const base64EncodedFile = im.split(",")[1];
        // console.log(base64EncodedFile); // base64にしたデータ
        this.images.push(base64EncodedFile);
        this.files.push(file);
      };
    },

    dropFile() {
      this.isEnter = false;
      const files = [...event.dataTransfer.files];
      this.checkTypeAndGPS(files);
    },
    onImageChange(e) {
      const files = e.target.files || e.dataTransfer.files;
      this.checkTypeAndGPS(files);
    },

    upload: function () {
      this.changeButton();
      var resStatus;
      var getRequest;

      axios
        .post(
          "/api/auth/trip/save",
          {
            user_id: this.$store.getters["auth/getUserID"],
            imgs: this.images,
          },
          {
            headers: {
              "X-Token": this.$store.getters["auth/getToken"],
            },
          }
        )
        .then((response) => {
          resStatus = response.status;
          getRequest = response.data.trip_id;
          // console.log(getRequest);
          this.error = "";
          if (resStatus == 200) {
            // console.log(this.$store.state.tripid);
            // console.log(getRequest);
            this.$store.commit("trip/setTripID", getRequest);
            // console.log(this.$store.state.tripid);
            this.$router.push("/map");
          }
        });
      // .catch((error) => {
      //   // console.log("error");
      //   // console.log(error);
      // });
    },
    onClickTitle: function () {
      this.$router.push("/home");
    },

    onClickOpenPopup: function () {
      this.popup = true;
    },

    onClickClosePopup: function () {
      this.popup = false;
    },
  },
};
</script>
<style scoped>
input {
  display: none;
}

img:hover {
  opacity: 0.7;
}

#error {
  color: red;
}
#howToUse {
  text-align: center;
  font-family: serif;
  color: #5c5c5c;
  margin-top: 10px;
}
label {
  display: flex;
  justify-content: center;
  align-items: center;
}
.drop_area {
  color: #5c5c5c;
  font-weight: bold;
  font-family: serif;
  font-size: 1.2em;
  display: flex;
  justify-content: center;
  align-items: center;
  /* width: 500px;
  height: 300px; */
  width: 80vw;
  height: 48vh;
  border: 5px solid #5c5c5c;
  border-radius: 15px;
}
.enter {
  border: 10px dotted #5c5c5c;
}
ul {
  margin: 0;
  padding: 0;
  list-style-type: none;
}
.flex {
  display: flex;
  align-items: center;
}

.flex-col {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 0.5em;
  font-size: 10px;
}

.delete-mark {
  position: absolute;
  top: -14px;
  right: -10px;
  font-size: 20px;
}
span {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.button {
  padding: 0.5em 1.5em;
  background-color: #0070a7;
  color: white;
  font-size: 14px;
  font-weight: bold;
  border-radius: 5px;
  border-color: #0070a7;
}

.header {
  height: 50px;
  width: 100%;
  background-color: #2d2d2d;
  background-size: cover;
  background-position: center center;
  display: flex;
  align-items: center;
}

.headerTitle {
  /* width: 170px; */
  margin-left: 10px;
  margin-right: 10px;
  font-family: serif;
  font-size: 30px;
  cursor: pointer;
  color: white;
  display: table-cell;
  vertical-align: middle;
  text-align: center;
}

.title {
  margin-left: 10px;
  margin-right: 10px;
}

.menu {
  display: flex;
  justify-content: space-around;
  align-items: center;
  margin-left: auto;
}

.toHome {
  font-family: serif;
  color: white;
  cursor: pointer;
  font-size: 1em;
  text-align: center;
}

.addtrip {
  font-family: serif;
  color: white;
  cursor: pointer;
  font-size: 20px;
  background-color: #52a7f4;
  padding: 8px;
  border-radius: 10px;
  margin-left: 20px;
  margin-right: 20px;
}

.userSetting {
  color: white;
  margin-left: 10px;
  margin-right: 20px;
  cursor: pointer;
  background-color: #c850bc;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  text-align: center;
  line-height: 40px;
  min-width: 40px;
}

.input_body {
  background-color: rgb(219, 215, 215);
  height: 100vh;
  display: flex;
  justify-content: space-around;
  align-items: center;
}

.input_area {
  background-color: rgb(241, 241, 241);
  height: 80vh;
  width: 90%;
  display: flex;
  flex-flow: column;
  justify-content: center;
  align-items: center;
}

@media screen and (max-width: 440px) {
  .header {
    height: 50px;
    width: 100%;
    background-color: #2d2d2d;
    background-size: cover;
    background-position: center center;
    display: flex;
    align-items: center;
  }

  .headerTitle {
    width: 100px;
    font-family: serif;
    font-size: 1.2em;
    cursor: pointer;
    color: white;
    display: table-cell;
    vertical-align: middle;
    text-align: center;
  }

  .userSetting {
    color: white;
    margin-left: 20px;
    margin-right: 20px;
    cursor: pointer;
    background-color: #c850bc;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    text-align: center;
    line-height: 40px;
    min-width: 40px;
  }

  .toHome {
    font-family: serif;
    color: white;
    cursor: pointer;
    font-size: 0.8em;
    text-align: center;
    margin-left: auto;
  }

  .menu {
    display: flex;
    justify-content: space-around;
    align-items: center;
    margin-left: auto;
  }

  .drop_area {
    color: #5c5c5c;
    font-weight: bold;
    font-family: serif;
    font-size: 1.2em;
    display: flex;
    justify-content: center;
    align-items: center;
    text-align: center;
    width: 80vw;
    height: 48vh;
    border: 5px solid #5c5c5c;
    border-radius: 15px;
  }
}
</style>
