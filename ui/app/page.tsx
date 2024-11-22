"use client";

import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import ConnectionForm from "@/components/ConnectionForm";
import { useState } from "react";
import axios, { AxiosError } from "axios";
import Constants from "@/lib/contants";

export default function Home() {
  const router = useRouter();
  const [sourceData, setSourceData] = useState(null);
  const [targetData, setTargetData] = useState(null);

  const handleFormData = (type: "source" | "target", data: any) => {
    if (type === "source") {
      setSourceData(data);
    } else {
      setTargetData(data);
    }
  };
  const handleConnect = async () => {
    if (!sourceData || !targetData) {
      toast.error("Please fill both source and target connection details");
      return;
    }

    try {
      const response = await axios.post(`${Constants.BE_API}/connect`,
        { source: sourceData, target: targetData },
        {
          headers: {
            "Content-Type": "application/json",
          }
        },
      );

      if (response?.status !== 200) {
        throw new Error("Connection failed");
      }

      toast.success("Successfully connected!");
      router.push("/dashboard");
    } catch (error) {
      if (error instanceof AxiosError) {
        error.response?.data?.message && toast.error(error.response.data.message);
      } else {
        toast.error("Failed to connect. Please try again.");
      }
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-5xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Connect Your Data</h1>
          <p className="mt-2 text-gray-600">
            Configure your source and target connections
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8 mb-8">
          <ConnectionForm
            title="Source Connection"
            type="source"
            onFormData={(data) => handleFormData("source", data)}
          />
          <ConnectionForm
            title="Target Connection"
            type="target"
            onFormData={(data) => handleFormData("target", data)}
          />
        </div>

        <div className="flex justify-center">
          <Button
            size="lg"
            className="px-8 py-6 text-lg"
            onClick={handleConnect}
          // disabled={!sourceData || !targetData}
          >
            Connect
          </Button>
        </div>
      </div>
    </div>
  );
}