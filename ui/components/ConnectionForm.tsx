"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import * as z from "zod";
import { Card } from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const formSchema = z.object({
  host: z.string().min(1, "Host is required"),
  port: z.string().regex(/^\d+$/, "Port must be a number"),
  user: z.string().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
  database: z.string().min(1, "Database is required"),
});

interface ConnectionFormProps {
  title: string;
  type: "source" | "target";
  onFormData: (data: z.infer<typeof formSchema>) => void;
}

export default function ConnectionForm({
  title,
  type,
  onFormData,
}: ConnectionFormProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      host: "",
      port: "",
      user: "",
      password: "",
      database: "",
    },
  });

  const handleFormChange = () => {
    const data = form.getValues();
    if (Object.values(data).every(value => value !== "")) {
      onFormData(data);
    }
  };

  return (
    <Card className="p-6">
      <h2 className="text-xl font-semibold mb-4">{title}</h2>
      <Form {...form}>
        <form onChange={handleFormChange} className="space-y-4">
          <FormField
            control={form.control}
            name="host"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Host</FormLabel>
                <FormControl>
                  <Input placeholder="localhost" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="port"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Port</FormLabel>
                <FormControl>
                  <Input placeholder="3306" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="user"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Username</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password</FormLabel>
                <FormControl>
                  <Input type="password" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="database"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Database Name</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </form>
      </Form>
    </Card>
  );
}