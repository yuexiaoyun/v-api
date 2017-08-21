package models

import (
	"fmt"
	"github.com/adam-hanna/arrayOperations"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ReturnVideoInfo []VideoInfo

func init() {

}

type VideoInfo struct {
	UserId            string            `json:"user_id"`
	UserAvatar        string            `json:"user_avatar"`
	UserNickname      string            `json:"user_nickname"`
	UserHomepage      string            `json:"user_homepage"`
	Vid               int64             `json:"vid"`
	VideoTitle        string            `json:"video_title"`
	VideoSubtitle     string            `json:"video_subtitle"`
	VideoCover        string            `json:"video_cover"`
	VideoCover375x375 string            `json:"video_cover_375_375"`
	VideoBigCover     string            `json:"video_big_cover"`
	VideoPlayNum      int64             `json:"video_play_num"`
	VideoCommentNum   int64             `json:"video_comment_num"`
	VideoDuration     string            `json:"video_duration"`
	VideoUrl          string            `json:"video_url"`
	VideoUploadTime   string            `json:"video_upload_time"`
	VideoChannel      string            `json:"video_channel"`
	VideoTags         string            `json:"video_tags"`
	VideoDefinitions  []VideoDefinition `json:"video_definitions"`
	VideoCategory     string            `json:"video_category"`
	VidCmsTime        int64             `json:"video_cms_time"`
}

type RetUserInfo struct {
	user_id       string
	user_avatar   string
	user_nickname string
	user_homepage string
}

func getTitle(rawVideo RawVideoInfo) string {
	var title string
	title = rawVideo.VideoTitle
	return title
}

func getVideoChannel(rawVideo RawVideoInfo) string {
	channel := rawVideo.Channel
	if channel == "" {
		channel = "vhuyaunknown"
	}
	return channel
}

//TODO 年周的计算方法有点繁琐，找时间改进
func getVideoCover(rawVideo RawVideoInfo) string {
	var cover string
	if len(rawVideo.Cover) == 0 {
		yearStr, weekStr := getYearAndWeek(rawVideo.UploadStartTime)
		cover = beego.AppConfig.String("videoCoverPotocal") + beego.AppConfig.String("videoCoverDomain") + "/" + yearStr + weekStr + "/" + strconv.FormatInt(rawVideo.Vid, 10) + "/4-220x124.jpg"
	} else {
		cover = rawVideo.Cover
	}
	return cover
}

//TODO 年周的计算方法有点繁琐，找时间改进
func getVideoCoverBySize(rawVideo RawVideoInfo, width string, height string) string {
	var cover string
	yearStr, weekStr := getYearAndWeek(rawVideo.UploadStartTime)
	cover = beego.AppConfig.String("videoCoverPotocal") + beego.AppConfig.String("videoCoverDomain") + "/" + yearStr + weekStr + "/" + strconv.FormatInt(rawVideo.Vid, 10) + "/4-" + width + "x" + height + ".jpg"
	return cover
}

//TODO 年周的计算方法有点繁琐，找时间改进
func getVideoBigCover(rawVideo RawVideoInfo) string {
	var cover string
	yearStr, weekStr := getYearAndWeek(rawVideo.UploadStartTime)
	cover = beego.AppConfig.String("videoCoverPotocal") + beego.AppConfig.String("videoCoverDomain") + "/" + yearStr + weekStr + "/" + strconv.FormatInt(rawVideo.Vid, 10) + "/4-640x360.jpg"
	return cover
}

func getYearAndWeek(timestamp int64) (string, string) {
	format := "2006-01-02 15:04:05"
	t, _ := time.Parse(format, time.Unix(timestamp, 0).Format(format))
	year, week := t.ISOWeek()
	year = year - 2000
	return fmt.Sprintf("%02d", year), fmt.Sprintf("%02d", week)
}

func getVideoInfo(rawVideo RawVideoInfo) VideoInfo {
	rawUser, status := GetRawUser(rawVideo.Yyuid)
	var retUserInfo RetUserInfo = RetUserInfo{}
	if status == "ok" {
		retUserInfo = RetUserInfo{
			user_id:       rawUser.user_id,
			user_avatar:   rawUser.user_avatar,
			user_nickname: rawUser.user_nickname,
			user_homepage: rawUser.user_homepage,
		}
	}

	videoInfo := VideoInfo{
		Vid:               rawVideo.Vid,
		VideoTitle:        getTitle(rawVideo),
		VideoSubtitle:     getTitle(rawVideo),
		VideoCover:        getVideoCover(rawVideo),
		VideoPlayNum:      rawVideo.VideoPlayNum,
		VideoCommentNum:   rawVideo.VideoSupport,
		VideoDuration:     getDuration(rawVideo),
		VideoBigCover:     getVideoBigCover(rawVideo),
		VideoCover375x375: getVideoCoverBySize(rawVideo, "375", "375"),
		VideoUrl:          getVideoUrl(rawVideo),
		VideoUploadTime:   getVideoUploadTime(rawVideo),
		VideoChannel:      getVideoChannel(rawVideo),
		VideoTags:         getVideoTags(rawVideo),
		VideoCategory:     getVideoCategory(rawVideo),
		VidCmsTime:        rawVideo.UploadStartTime,
		UserId:            retUserInfo.user_id,
		UserAvatar:        retUserInfo.user_avatar,
		UserNickname:      retUserInfo.user_nickname,
		UserHomepage:      retUserInfo.user_homepage,
	}
	videoDefinitions, status := GetVideoDefinitions(rawVideo.Vid, false, "1000,1300,350,yuanhua")
	if status == "ok" {
		videoInfo.VideoDefinitions = videoDefinitions
	} else {
		videoInfo.VideoDefinitions = []VideoDefinition{}
	}
	return videoInfo
}

func getVideoUrl(rawVideo RawVideoInfo) string {
	return beego.AppConfig.String("baseUrl") + "/play/" + strconv.FormatInt(rawVideo.Vid, 10) + ".html"
}

func getVideoUploadTime(rawVideo RawVideoInfo) string {
	return fmt.Sprint(time.Unix(rawVideo.UploadStartTime, 0).Format("2006-01-02 03:04:05"))
}
func getVideoTags(rawVideo RawVideoInfo) string {
	videoTags := GetTagByVid(rawVideo.Vid)
	var tags []string
	for _, tag := range videoTags {
		tags = append(tags, tag.Tag)
	}
	if len(tags) == 0 {
		return ""
	} else {
		return strings.Join(tags, ",")
	}
}

func getVideoPlayNum(rawVideo RawVideoInfo) int {
	return 0
}
func getUserInfo(uid int64) RetUserInfo {
	rawUser, status := GetRawUser(uid)
	var retUserInfo RetUserInfo = RetUserInfo{}
	if status == "ok" {
		retUserInfo = RetUserInfo{
			user_id:       rawUser.user_id,
			user_avatar:   rawUser.user_avatar,
			user_nickname: rawUser.user_nickname,
			user_homepage: rawUser.user_homepage,
		}
		return retUserInfo
	}
	return retUserInfo
}

func getDuration(rawVideo RawVideoInfo) string {
	input := rawVideo.Duration
	if input == "0" {
		return "00:00"
	}
	input = input + "s"
	duration, err := time.ParseDuration(input)
	if err != nil {
		beego.Info("转换时长出错")
		beego.Info(err)
		return "00:00"
	} else {
		durationFloat, err := strconv.ParseFloat(input, 64)
		if err != nil {
			var seconds int
			var minutes int
			var hours int
			if durationFloat >= 3600 {
				seconds = int(duration.Seconds()) % 60
				minutes = int(duration.Minutes()) % 60
				hours = int(duration.Hours()) % 24
				return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

			} else {
				seconds = int(duration.Seconds()) % 60
				minutes = int(duration.Minutes()) % 60
				return fmt.Sprintf("%02d:%02d", minutes, seconds)
			}
		} else {
			return "00:00"
		}
	}
}

func getVideoCategory(rawVideo RawVideoInfo) string {
	return GetVideoCategory(rawVideo.Channel)
}

func GetList(vidsList []int, limit int) []VideoInfo {
	var videoList []VideoInfo
	/*videoSize := len(vidsList)

	wg := sync.WaitGroup{}

	ReturnVideoInfo = []VideoInfo{}
	for i := 0; i < videoSize; i++ {
		wg.Add(1)
		go GoGetVideoSlice(&wg, vidsList[i])
		if limit != 0 && len(ReturnVideoInfo) >= limit {
			break
		}
	}
	wg.Wait()
	videoList = ReturnVideoInfo*/
	rawVideoInfo := GetRawVideoByList(vidsList)
	/*for _, vid := range vidsList {
		vidStr := strconv.Itoa(vid)
		videoInfo := GetByVid(vidStr)
		videoList = append(videoList, videoInfo)
		if limit != 0 && len(videoList) >= limit {
			break
		}
	}*/
	wg := sync.WaitGroup{}
	for _, rawVideoInfo := range rawVideoInfo {
		wg.Add(1)
		go GetByRawVideoInfo(rawVideoInfo, &wg, videoList)
		/*videoList = append(videoList, videoInfo)
		if limit != 0 && len(videoList) >= limit {
			break
		}*/
	}
	wg.Wait()
	return videoList
}

func GetVideoByUid(yyuid string, limit int, page int) []VideoInfo {
	cacheKey := SHENJTLIVE
	cacheKey = cacheKey + yyuid
	cacheHandler, errMsg := GetCacheHandler()
	var videoInfoList []VideoInfo
	if errMsg == nil {
		if _, _, e := cacheHandler.Get(cacheKey, &videoInfoList); e != nil {
			videoInfoList = GetVideoByUidFromDB(yyuid, limit, page)
			beego.Info("[GetVideoByUid]数据从表读取：")
			beego.Info(videoInfoList)
			//判断结构vid是否为空，不空，设置缓存
			if len(videoInfoList) != 0 {
				SetDataIntoCache(cacheHandler, cacheKey, videoInfoList, SHENJTLIVE_TIMEOUT)
			}
		} else {
			beego.Info("[GetVideoByUid]数据从缓存读取：")
			beego.Info(videoInfoList)
		}

	} else {
		beego.Error("[GetVideoByUid]获取缓存句柄失败")
		videoInfoList = GetVideoByUidFromDB(yyuid, limit, page)
	}
	return videoInfoList

}

func GetVideoByUidFromDB(yyuid string, limit int, page int) []VideoInfo {
	liveVids := GetDotVidByUid(yyuid, limit, page)
	uploadVids := GetVidsByUid(yyuid, limit, page)
	var vidsIntList []int
	vidsIntList = mergeAndSort(liveVids, uploadVids)
	videoInfoList := GetList(vidsIntList, 20)
	return videoInfoList
}

func mergeAndSort(liveVids []int, uploadVids []int) []int {
	var vidsIntList []int
	if len(liveVids) != 0 && len(uploadVids) != 0 {
		var vids reflect.Value
		var ret bool
		vids, ret = arrayOperations.Union(liveVids, uploadVids)
		if ret {
			vidsIntList, ret = vids.Interface().([]int)
			if !ret {
				beego.Error("vid数组转化失败")
				return []int{}
			} else {
				//倒序
				sort.Sort(sort.Reverse(sort.IntSlice(vidsIntList)))
				return vidsIntList
			}
		} else {
			beego.Error("vid数组合并失败")
			return []int{}
		}
	} else {
		if len(uploadVids) != 0 {
			sort.Sort(sort.Reverse(sort.IntSlice(uploadVids)))
			vidsIntList = make([]int, len(uploadVids))
			copy(vidsIntList, uploadVids[:])
			return vidsIntList
		} else if len(liveVids) != 0 {
			sort.Sort(sort.Reverse(sort.IntSlice(liveVids)))
			vidsIntList = make([]int, len(liveVids))
			copy(vidsIntList, liveVids[:])
			return vidsIntList
		} else {
			return []int{}
		}
	}
}

//异步获取，后面发现可能不符合需求
func GoGetVideoSlice(wg *sync.WaitGroup, vid int) {
	videoInfo := GetByVid(strconv.Itoa(vid))
	ReturnVideoInfo = append(ReturnVideoInfo, videoInfo)
	wg.Done()
}

func GetByVid(vid string) VideoInfo {
	cacheKey := VIDEOINFO
	cacheKey = cacheKey + vid
	cacheHandler, errMsg := GetCacheHandler()
	var videoInfo VideoInfo
	rawVideo := GetRawVideo(vid)
	if errMsg != nil {
		videoInfo = getVideoInfo(rawVideo)
		beego.Info("数据从表读取：")
		beego.Info(videoInfo)
		//判断结构vid是否为空，不空，设置缓存
		if videoInfo.Vid != 0 {
			SetDataIntoCache(cacheHandler, cacheKey, videoInfo, VIDEOINFO_TIMEOUT)
		}
	} else {
		if _, _, e := cacheHandler.Get(cacheKey, &videoInfo); e != nil {
			beego.Info("解析有问题")
			beego.Info(e)
			videoInfo = getVideoInfo(rawVideo)
			beego.Info("数据从表读取：")
			beego.Info(videoInfo)
			//判断结构vid是否为空，不空，设置缓存
			if videoInfo.Vid != 0 {
				SetDataIntoCache(cacheHandler, cacheKey, videoInfo, VIDEOINFO_TIMEOUT)
			}
		} else {
			beego.Info("数据从缓存读取：")
			beego.Info(videoInfo)
		}
	}

	return videoInfo
}

func GetByRawVideoInfo(rawVideo RawVideoInfo, wg *sync.WaitGroup, videoList []VideoInfo) {
	cacheKey := VIDEOINFO
	cacheKey = cacheKey + strconv.Itoa(int(rawVideo.Vid))
	cacheHandler, errMsg := GetCacheHandler()
	var videoInfo VideoInfo
	if errMsg != nil {
		videoInfo = getVideoInfo(rawVideo)
		beego.Info("数据从表读取：")
		beego.Info(videoInfo)
		//判断结构vid是否为空，不空，设置缓存
		if videoInfo.Vid != 0 {
			SetDataIntoCache(cacheHandler, cacheKey, videoInfo, VIDEOINFO_TIMEOUT)
		}
	} else {
		if _, _, e := cacheHandler.Get(cacheKey, &videoInfo); e != nil {
			beego.Info("解析有问题")
			beego.Info(e)
			videoInfo = getVideoInfo(rawVideo)
			beego.Info("数据从表读取：")
			beego.Info(videoInfo)
			//判断结构vid是否为空，不空，设置缓存
			if videoInfo.Vid != 0 {
				SetDataIntoCache(cacheHandler, cacheKey, videoInfo, VIDEOINFO_TIMEOUT)
			}
		} else {
			beego.Info("数据从缓存读取：")
			beego.Info(videoInfo)
		}
	}
	videoList = append(videoList, videoInfo)
}
