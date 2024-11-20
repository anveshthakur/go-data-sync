"use client";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import Loader from "@/components/ui/loader";
import { useEffect, useState } from "react";
import DataList from "./DataList";
import DataTable from "./DataTable";
import axios from "axios"
import { toast } from "sonner";

export default function TabsSection() {
  const [selectedItem, setSelectedItem] = useState('');
  const [loader, setLoader] = useState(true);
  const [tableList, setTableList] = useState([]);

  useEffect(() => {
    (async () => {
      try {
        const response = await axios.get('http://localhost:8080/tables?type=source');
        setTableList(response?.data);
        setSelectedItem(response?.data[0]);
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
              >
                Source
              </TabsTrigger>
            </TabsList>

            <TabsContent value="source" className="flex-grow mt-0 p-0">
              <DataList
                items={tableList}
                onItemClick={handleItemClick}
                selectedItem={selectedItem}
              />
            </TabsContent>

          </Tabs>

          {selectedItem && (
            <DataTable item={selectedItem} setLoader={setLoader} />
          )}
        </div>
      )
      }
    </>

  );
}