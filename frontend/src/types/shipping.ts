export interface CourierRate {
  courierName: string
  serviceCode: string
  serviceName: string
  price: number
  deliveryEta: string
}

export interface ShipmentInfo {
  courierName: string
  shippingRate: number
  status: string
  trackingNumber: string
}

export interface RateCheckParams {
  receiverPostcode: string
  weightKg: number
}
