import { Bus, Truck, Crown } from "lucide-react";
import { Car } from "lucide-react";
import { CarPackageSlug } from "../types";

export const PackagesMeta: Record<CarPackageSlug, {
  name: string,
  icon: React.ReactNode,
  description: string,
}> = {
  [CarPackageSlug.BIKE]: {
    name: "Bike",
    icon: <Car />,
    description: "Affordable solo rides",
  },
  [CarPackageSlug.AUTO]: {
    name: "Auto",
    icon: <Truck />,
    description: "Affordable rides for 2-3 people",
  },
  [CarPackageSlug.SEDAN]: {
    name: "Sedan",
    icon: <Bus />,
    description: "Comfortable rides for up to 4 people",
  },
  [CarPackageSlug.SUV]: {
    name: "SUV",
    icon: <Crown />,
    description: "Spacious rides for up to 6 people",
  },
}