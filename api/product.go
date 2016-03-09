package api

import (
    "fmt"
    "net/url"
    "strconv"
    "github.com/mitchellh/mapstructure"
)

const (
    DEFAULT_PRODUCT_API_LENGTH = 20
    DEFAULT_PRODUCT_MAX_LENGTH = 100
    
    DEFAULT_PRODUCT_MAX_OFFSET = 50000
)

type ProductService struct {
    ApiId        string `mapstructure:"api_id"`
    AffiliateId  string `mapstructure:"affiliate_id"`
    Site         string `mapstructure:"site"`
    Service      string `mapstructure:"service"`
    Floor        string `mapstructure:"floor"`
    Length       int64  `mapstructure:"hits"`
    Offset       int64  `mapstructure:"offset"`
    Sort         string `mapstructure:"sort"`
    Keyword      string `mapstructure:"keyword"`
}

type ProductRawResponse struct {
    Request ProductService  `mapstructure:"request"`
    Result  ProductResponse `mapstructure:"result"`
}

type ProductResponse struct {
    ResultCount   int64  `mapstructure:"result_count"`
    TotalCount    int64  `mapstructure:"total_count"`
    FirstPosition int64  `mapstructure:"first_position"`
    Items         []Item `mapstructure:"items"`
}

type Item struct {
    AffiliateUrl       string             `mapstructure:"affiliateURL"`
    AffiliateUrlMobile string             `mapstructure:"affiliateURLsp"`
    CategoryName       string             `mapstructure:"category_name"`
    Comment            string             `mapstructure:"comment"`
    ContentId          string             `mapstructure:"content_id"`
    Date               string             `mapstructure:"date"`
    FloorName          string             `mapstructure:"floor_name"`
    ISBN               string             `mapstructure:"isbn"`
    JANCode            string             `mapstructure:"jancode"`
    ProductCode        string             `mapstructure:"maker_product"`
    ProductId          string             `mapstructure:"product_id"`
    ServiceName        string             `mapstructure:"service_name"`
    Stock              string             `mapstructure:"stock"`
    Title              string             `mapstructure:"title"`
    Url                string             `mapstructure:"URL"`
    UrlMoble           string             `mapstructure:"URLsp"`
    Volume             string             `mapstructure:"volume"`
    ImageUrl           ImageUrlList       `mapstructure:"imageURL"`
    SampleImageUrl     SampleImageUrlList `mapstructure:"sampleImageURL"`
    SampleMovieUrl     SampleMovieUrlList `mapstructure:"sampleMovieURL"`
    Review             ReviewInformation  `mapstructure:"review"`
    PriceInformation   PriceInformation   `mapstructure:"prices"`
    ItemInformation    ItemInformation    `mapstructure:"iteminfo"`
    BandaiInformation  BandaiInformation  `mapstructure:"bandaiinfo"`
    CdInformation      CdInformation      `mapstructure:"cdinfo"`
}

type ImageUrlList struct {
    List  string `mapstructure:"list"`
    Small string `mapstructure:"small"`
    Large string `mapstructure:"large"`
}

type SampleImageUrlList struct {
    Sample_s SmallSampleList  `mapstructure:"sample_s"`
}

type SmallSampleList struct {
    Image []string `mapstructure:"image"`
}

type SampleMovieUrlList struct {
    Size_476_306 string `mapstructure:"size_476_306"`
    Size_560_360 string `mapstructure:"size_560_360"`
    Size_644_414 string `mapstructure:"size_644_414"`
    Size_720_480 string `mapstructure:"size_720_480"`
    PC_flag      bool   `mapstructure:"pc_flag"`
    SP_flag      bool   `mapstructure:"sp_flag"`
}

type PriceInformation struct {
    Price       string             `mapstructure:"price"`
    PriceAll    string             `mapstructure:"price_all"`
    RetailPrice string             `mapstructure:"list_price"`
    Distributions DistributionList `mapstructure:"deliveries"`
}

type DistributionList struct {
    Distribution []Distribution `mapstructure:"delivery"`
}

type Distribution struct {
    Type  string `mapstructure:"type"`
    Price string `mapstructure:"price"`
}

type ItemInformation struct {
    Maker     []ItemComponent `mapstructure:"maker"`
    Label     []ItemComponent `mapstructure:"label"`
    Series    []ItemComponent `mapstructure:"series"`
    Keywords  []ItemComponent `mapstructure:"keyword"`
    Genres    []ItemComponent `mapstructure:"genre"`
    Actors    []ItemComponent `mapstructure:"actor"`
    Artists   []ItemComponent `mapstructure:"artist"`
    Authors   []ItemComponent `mapstructure:"author"`
    Directors []ItemComponent `mapstructure:"director"`
    Fighters  []ItemComponent `mapstructure:"fighter"`
    Colors    []ItemComponent `mapstructure:"color"`
    Sizes     []ItemComponent `mapstructure:"size"`
}

type ItemComponent struct {
    Id   string `mapstructure:"id"`
    Name string `mapstructure:"name"`
}

type BandaiInformation struct {
    TitleCode string `mapstructure:"titlecode"`
}

type CdInformation struct {
    Kind string `mapstructure:"kind"`
}

type ReviewInformation struct {
    Count   int64   `mapstructure:"count"`
    Average float64 `mapstructure:"average"`
}

func NewProductService(affiliateId, apiId string) *ProductService {
    return &ProductService{
        ApiId:       apiId,
        AffiliateId: affiliateId,
        Site:        "",
        Service:     "",
        Floor:       "",
        Length:      DEFAULT_PRODUCT_API_LENGTH,
        Offset:      DEFAULT_API_OFFSET,
        Sort:        "",
        Keyword:     "",
    }
}

func (srv *ProductService) Execute() (*ProductResponse, error) {
    result, err := srv.ExecuteWeak()
    if err != nil {
        return nil, err
    }
    var raw ProductRawResponse
    if err = mapstructure.WeakDecode(result, &raw); err != nil {
        return nil, err
    }
    return &raw.Result, nil
}

func (srv *ProductService) ExecuteWeak() (interface{}, error) {
    reqUrl, err := srv.BuildRequestUrl()
    if err != nil {
        return nil, err
    }

    return RequestJson(reqUrl)
}

func (srv *ProductService) SetLength(length int64) *ProductService {
    srv.Length = length
    return srv
}

func (srv *ProductService) SetHits(length int64) *ProductService {
    srv.SetLength(length)
    return srv
}

func (srv *ProductService) SetOffset(offset int64) *ProductService {
    srv.Offset = offset
    return srv
}

func (srv *ProductService) SetKeyword(keyword string) *ProductService {
    srv.Keyword = TrimString(keyword)
    return srv
}

func (srv *ProductService) SetSort(sort string) *ProductService {
    srv.Sort = TrimString(sort)
    return srv
}

func (srv *ProductService) SetSite(site string) *ProductService {
    srv.Site = TrimString(site)
    return srv
}

func (srv *ProductService) SetService(service string) *ProductService {
    srv.Service = TrimString(service)
    return srv
}

func (srv *ProductService) SetFloor(floor string) *ProductService {
    srv.Floor = TrimString(floor)
    return srv
}

func (srv *ProductService) ValidateLength() bool {
    return ValidateRange(srv.Length, 1, DEFAULT_PRODUCT_MAX_LENGTH)
}

func (srv *ProductService) ValidateOffset() bool {
    return ValidateRange(srv.Offset, 1, DEFAULT_PRODUCT_MAX_OFFSET)
}

func (srv *ProductService) BuildRequestUrl() (string, error) {
    if srv.ApiId == "" {
        return "", fmt.Errorf("set invalid ApiId parameter.")
    }
    if !ValidateAffiliateId(srv.AffiliateId) {
        return "", fmt.Errorf("set invalid AffiliateId parameter.")
    }

    if !ValidateSite(srv.Site) {
        return "", fmt.Errorf("set invalid Site parameter.")
    }

    queries := url.Values{}
    queries.Set("api_id", srv.ApiId)
    queries.Set("affiliate_id", srv.AffiliateId)
    queries.Set("site", srv.Site)

    if srv.Length != 0 {
        if !srv.ValidateLength() {
            return "", fmt.Errorf("length out of range: %d", srv.Length)
        }
        queries.Set("hits", strconv.FormatInt(srv.Length, 10))
    }

    if srv.Offset != 0 {
        if !srv.ValidateOffset() {
            return "", fmt.Errorf("offset out of range: %d", srv.Offset)
        }
        queries.Set("offset", strconv.FormatInt(srv.Offset, 10))
    }

    if (srv.Service != "") {
        queries.Set("service", srv.Service)
    }
    if (srv.Floor != "") {
        queries.Set("floor", srv.Floor)
    }
    if (srv.Sort != "") {
        queries.Set("sort", srv.Sort)
    }
    if (srv.Keyword != "") {
        queries.Set("keyword", srv.Keyword)
    }
    return API_BASE_URL + "/ItemList?" + queries.Encode(), nil
}