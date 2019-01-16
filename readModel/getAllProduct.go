package readModel
import (
    "fmt"
	"github.com/go-pg/pg"
	"time"
	"github.com/golang/protobuf/proto"
	"proto_pg_orm/protoModel"
)

// Initialize Result Data for postgresql ORM (Products)
type Result struct {
	Id int64
	Name string
	Description string
	ListPrice float64 
	Cost float64
	Volume float64 
	Weight float64
	CategoryId int64
	Category string
	CategoryCode string
	UomSalesId int64
	UomSales string
	UomPurchaseId int64
	UomPurchase string
	Sku string
	BrandId int64
	Brand string
	IsBundle bool
	CreateDate time.Time
}
// Initialize Result Data for postgresql ORM Bundle Products
type ResultBundle struct {
	Id int64
	Name string
	Description string
	ListPrice float64
	Cost float64
}

// Query to get all data
func GetResult(db *pg.DB, query string) ([]Result, error) {
    var result []Result
    _, err := db.Query(&result, query)
    return result, err
}

// Query to get all data by Id
func GetResultByIds(db *pg.DB, query string, ids []int64) ([]ResultBundle, error) {
    var resultBundle []ResultBundle
    _, err := db.Query(&resultBundle, query, pg.In(ids))
    return resultBundle, err
}

// Gathering Data
func GetAllProduct() (dataByte []byte){
	// Declaring Protobuf Message
	apR := &protoModel.AllProduct{}
	// Connecting to Postgresql
	db := pg.Connect(&pg.Options{
		Addr: "localhost:5432",
		User: "henis",
		Password: "12345678",
		Database: "oktagon",
	})
	defer db.Close()
	// Getting the result data
	results, err := GetResult(db, `select distinct on (pt.id) pt.id, pt.name, cast(pt.description as varchar), pt.list_price, pph.cost, pt.volume, pt.weight, pt.categ_id as category_id, pc.name as category, pc.code as category_code, pt.uom_id as uom_sales_id, pus.name as uom_sales, pt.uom_po_id as uom_purchase_id, pup.name uom_purchase, pt.default_code as sku, pt.product_brand_id as brand_id, rp.name as brand, pt.is_bundle, pph.create_date from product_template pt left outer join product_price_history pph on pt.id = pph.product_id inner join product_category pc on pt.categ_id = pc.id inner join product_uom pus on pt.uom_id = pus.id inner join product_uom pup on pt.uom_po_id = pup.id inner join product_brand pb on pt.product_brand_id = pb.id inner join res_partner rp on pb.partner_id = rp.id order by pt.id asc, pph.create_date desc`)
	if err != nil{
		fmt.Println(err)
	}
	// Implementing result data to protobuf
	for _, element := range results{
		pR := &protoModel.Product{
			Id: element.Id,
			Name: element.Name,
			Description: element.Description,
			ListPrice: element.ListPrice,
			Cost: element.Cost,	
			Volume: element.Volume,
			Weight: element.Weight,
			CategoryId: element.CategoryId,
			Category: element.Category,
			CategoryCode: element.CategoryCode,
			UomSalesId: element.UomSalesId,
			UomSales: element.UomSales,
			UomPurchaseId: element.UomPurchaseId,
			UomPurchase: element.UomPurchase,
			Sku: element.Sku,
			BrandId: element.BrandId,
			Brand: element.Brand,
			IsBundle: element.IsBundle,
			CreateDate: element.CreateDate.Unix(),
		}
		// Gathering data for bundling product
		if element.IsBundle == true{
			var elementId []int64
			elementId = append(elementId, element.Id)
			resultsBundle, err := GetResultByIds(db, `select pt.id, pt.name, pt.description, pt.list_price, pph.cost from product_template pt left outer join product_price_history pph on pt.id = pph.product_id left outer join product_template_bundle ptb on pt.id = ptb.product_bundle_ids where ptb.products_in_bundle_ids = (?)`, elementId)
			if err != nil{
				fmt.Println(err)
			}
			for  _, elementbundle := range resultsBundle{
				pBR := &protoModel.ProductBundleDetail{
					Id: elementbundle.Id,
					Name: elementbundle.Name,
					Description: elementbundle.Description,
					ListPrice: elementbundle.ListPrice,
					Cost: elementbundle.Cost,
				}
				pR.ProductBundleDetail = append(pR.ProductBundleDetail, pBR)
			}
		}
		apR.Product = append(apR.Product, pR)
	}
	// Marshaling protobuf message
	data, err := proto.Marshal(apR)
	if err != nil {
		fmt.Println(err)
	}
	// Returning protobuf message
	dataByte = data
	return
}