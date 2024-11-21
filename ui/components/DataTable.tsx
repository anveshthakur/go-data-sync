"use client";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import { Button } from "@/components/ui/button";
import axios from "axios";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import Constants from "@/lib/contants";

export default function DataTable({ item, setLoader }: { item: any, setLoader: (loader: boolean) => void }) {
  const [tableDetails, setTableDetails] = useState([]);

  useEffect(() => {
    (async () => {
      const response = await axios.get(`${Constants.BE_API}/fetch-data?table=${item}&type=source`);
      setTableDetails(response?.data);
      setLoader(false);
    })();
  }, [item]);

  const synchronizeData = async () => {
    try {
      const response = await axios.post(`${Constants.BE_API}/sync`,
        [
          { type: "source", "table": item },
          { type: "target", "table": item }
        ],
        {
          headers: {
            "Content-Type": "application/json",
          }
        },
      );

      if (response?.status !== 201) {
        throw new Error("Synchronization failed");
      }

      toast.success("Synchronization successful!")
    } catch (error) {
      toast.error("Failed to synchronize. Please try again.");
    }
  };

  return (
    <div className="absolute top-0 left-[20%] right-0 h-full bg-white p-6 animate-in slide-in-from-right duration-300">
      <div className="max-w-5xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-semibold">{item}</h2>
          <div className="flex justify-end">
            <Button
              size="lg"
              className="px-8 py-6 text-lg"
              onClick={synchronizeData}
              disabled={tableDetails?.length === 0 || !tableDetails}
            >
              Synchronize Table
            </Button>
          </div>

        </div>
        <div className="border rounded-lg overflow-hidden">

          {tableDetails?.length > 0 &&
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    {
                      Object.keys(tableDetails[0]).map((columnName, id) => <TableHead key={id}>{columnName}</TableHead>)
                    }
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {tableDetails.map((detail, id) => (
                    <TableRow key={id}>
                      {
                        Object.values(detail).map((value, id) => <TableCell key={id}>{value as string}</TableCell>)
                      }
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </>
          }
        </div>
      </div>
    </div>
  );
}