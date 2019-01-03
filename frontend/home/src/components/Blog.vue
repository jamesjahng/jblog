<template>
  <div>
    <div v-if="year !== undefined" v-html="articleHtmlSource">
    </div>
    <div v-else>
      <div class="blog-list-element" v-for="i in index" :key="i">
        <router-link :to="i.uri">{{ i.title }}</router-link>
        <br>
        {{ i.date }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'blogExample',
  data() {
    return {
      // __INSERTION_POSITION__ // DONT CHANGE!!
      index: 
[{"title":"[웹] 블로깅 시스템을 만들었습니다","uri":"/blog/2019/01/02/start-blog/","date":"2019/01/02"}] // __INSERTION_POSITION_END__ // DONT CHANGE!!
,
    year : this.$route.params.year,
    month : this.$route.params.month,
    day : this.$route.params.day,
    title : this.$route.params.title,
    articleHtmlSource : ""
    }
  },
    watch: {
    '$route' (to) {
      this.year = to.params.year;
      this.month = to.params.month;
      this.day = to.params.day;
      this.title = to.params.title;
      if (this.year === undefined) {
        return;
      }
      var htmlDocUri = 
        '/blog_contents/'
        + this.year + '/'
        + this.month + '/'
        + this.day + '/'
        + this.title + '.html';
      fetch(htmlDocUri)
        .then(response => response.text())
        .then(responseText => this.articleHtmlSource = responseText);
    }
  }
}
</script>

<style scoped>
.blog-list-element {
  margin-bottom: 10px;
}
</style>
