"use client";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import Loader from "@/components/ui/loader";
import { useEffect, useState } from "react";
import DataList from "./DataList";
import DataTable from "./DataTable";
import axios from "axios"
import { toast } from "sonner";
import Link from "next/link";
import Constants from "@/lib/contants";

export default function TabsSection() {
  const [selectedItem, setSelectedItem] = useState('');
  const [type, setType] = useState('');
  const [loader, setLoader] = useState(true);
  const [tableList, setTableList] = useState([]);

  useEffect(() => {
    (async () => {
      try {
        const response = await axios.get(`${Constants.BE_API}/tables?type=source`);
        setTableList(response?.data);
        setSelectedItem(response?.data[0]);
        setType("source");
        setLoader(false);
      } catch (error) {
        setLoader(false);
        toast.error("Failed to fetch tables. Please try again.");
      }
    })();
  }, []);


  const handleItemClick = (item: string) => {
    setSelectedItem(item);
  };

  const handleTabClick = async (type: string) => {
    const response = await axios.get(`${Constants.BE_API}/tables?type=${type}`);
    setType(type);
    setTableList(response?.data);
    setSelectedItem(response?.data[0]);
  };

  return (
    <>
      {
        loader && <Loader />
      }
      {!loader && tableList.length > 0 && (

        <div className="h-full">
          <Tabs defaultValue="source" className="h-full flex flex-col">
            <TabsList className="w-full justify-start rounded-none border-b bg-transparent">
              <TabsTrigger
                value="source"
                className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
                onClick={() => handleTabClick("source")}
              >
                Source
              </TabsTrigger>
              <TabsTrigger
                value="target"
                className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
                onClick={() => handleTabClick("target")}
              >
                Target
              </TabsTrigger>
              <Link href="/"><h1 className="ml-28 cursor-pointer text-black font-semibold">Home</h1></Link>
            </TabsList>
            <TabsContent value="source" className="flex-grow mt-0 p-0">
              <DataList
                items={tableList}
                onItemClick={handleItemClick}
                selectedItem={selectedItem}
              />
            </TabsContent>
            <TabsContent value="target" className="flex-grow mt-0 p-0">
              <DataList
                items={tableList}
                onItemClick={handleItemClick}
                selectedItem={selectedItem}
              />
            </TabsContent>
          </Tabs>

          {selectedItem && (
            <DataTable item={selectedItem} type={type} setLoader={setLoader} />
          )}
        </div>
      )
      }
    </>

  );
}