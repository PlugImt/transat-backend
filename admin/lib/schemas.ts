import { z } from "zod";

export const userSchema = z.object({
  email: z.string().min(1, "L'email est requis").email("Email invalide"),
  first_name: z
    .string()
    .min(1, "Le prénom est requis")
    .min(2, "Le prénom doit contenir au moins 2 caractères"),
  last_name: z
    .string()
    .min(1, "Le nom est requis")
    .min(2, "Le nom doit contenir au moins 2 caractères"),
  phone_number: z
    .string()
    .optional()
    .refine((val) => !val || /^[+]?[\d\s\-()]+$/.test(val), {
      message: "Numéro de téléphone invalide",
    }),
  campus: z.string().optional(),
  formation_name: z.string().optional(),
  graduation_year: z
    .union([z.string(), z.number()])
    .optional()
    .transform((val) => {
      if (!val || val === "" || val === "0") return undefined;
      const num = typeof val === "string" ? parseInt(val, 10) : val;
      return num;
    })
    .refine(
      (val) => {
        if (val === undefined) return true;
        const currentYear = new Date().getFullYear();
        return val >= 1900 && val <= currentYear + 10;
      },
      {
        message: "Année de diplôme invalide",
      },
    ),
  language: z.string().optional(),
  roles: z.array(z.string()).default([]),
});

export type UserFormData = z.infer<typeof userSchema>;

export const eventSchema = z.object({
  name: z.string().min(1, "Le nom est requis"),
  description: z.string().min(1, "La description est requise"),
  location: z.string().min(1, "Le lieu est requis"),
  link: z.string().url("Lien invalide").optional().or(z.literal("")),
  start_date: z.string().min(1, "La date de début est requise"),
  end_date: z.string().min(1, "La date de fin est requise"),
  picture: z.string().optional(),
  id_club: z.number().optional(),
});

export type EventFormData = z.infer<typeof eventSchema>;

export const clubSchema = z.object({
  name: z.string().min(1, "Le nom est requis"),
  description: z.string().min(1, "La description est requise"),
  location: z.string().min(1, "Le lieu est requis"),
  link: z.string().url("Lien invalide").optional().or(z.literal("")),
  picture: z.string().optional(),
});

export type ClubFormData = z.infer<typeof clubSchema>;
